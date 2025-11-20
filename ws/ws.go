package ws

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yanun0323/errors"
	"github.com/yanun0323/logs"
	"github.com/yanun0323/pkg/channel"
	"github.com/yanun0323/pkg/sys"
)

const (
	_defaultMessageQueueCap int = 1_000
)

type WebSocket struct {
	conn atomic.Value // *dialing
	url  string

	dial      func() (*dialing, error)
	shutdown  chan struct{}
	reconnect chan struct{}

	registers []Sidecar

	subscribers     map[uint64]chan Message
	subscribersLock sync.RWMutex
	nextID          atomic.Uint64

	start  atomic.Bool
	end    atomic.Bool
	logger logs.Logger
}

type Sidecar struct {
	Sender  func(context.Context, *WebSocket) error
	Waiter  func(context.Context, Message) (isExpected bool, failure error)
	Timeout time.Duration
}

func New(ctx context.Context, url string, ping ...bool) *WebSocket {
	ws := &WebSocket{
		url: url,
		dial: func() (*dialing, error) {
			return dial(ctx, url, ping...)
		},
		shutdown:    make(chan struct{}),
		reconnect:   make(chan struct{}, 1),
		subscribers: make(map[uint64]chan Message),
		logger: logs.Get(ctx).With(
			"websocket", newLogID(),
			"url", url,
		),
	}

	return ws
}

func ReadMessage[T any](msg Message) (T, bool) {
	var resp T
	err := json.Unmarshal(msg.Data, &resp)
	if err != nil {
		if errors.As(err, json.UnmarshalTypeError{}) {
			logs.Debugf("unmarshal message: mismatch json type, err: %+v", err)
		} else {
			logs.Debugf("unmarshal message, err: %+v", err)
		}
	}

	return resp, err == nil
}

const (
	DefaultTimeout = 15 * time.Second
)

func (ws *WebSocket) SendAndWait(ctx context.Context, executor Sidecar) error {
	if executor.Sender == nil || executor.Waiter == nil {
		return errors.New("invalid hook, require sender and waiter")
	}

	done := make(chan error, 1)
	defer channel.SafeClose(done)

	waitTimeout := executor.Timeout
	if executor.Timeout <= 0 {
		waitTimeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()

	msgCh, unsubscribe := ws.Subscribe()
	defer unsubscribe()

	go func() {
		for {
			select {
			case msg, ok := <-msgCh:
				if !ok {
					channel.TryPush(done, error(errors.New("message channel closed")))
					return
				}

				ok, err := executor.Waiter(ctx, msg)
				if err != nil {
					channel.TryPush(done, error(errors.Wrap(err, "waiting for message")))
					return
				}

				if ok {
					channel.TryPush(done, nil)
					return
				}
			case <-ctx.Done():
				channel.TryPush(done, ctx.Err())
				return
			case <-sys.Shutdown():
				channel.TryPush(done, context.Canceled)
				return
			}
		}
	}()

	if err := executor.Sender(ctx, ws); err != nil {
		return errors.Wrap(err, "send func error before waiting")
	}

	return errors.Wrap(<-done, "websocket message")
}

func (ws *WebSocket) clearConnection() {
	var d *dialing
	ws.conn.Store(d)
}

func (ws *WebSocket) observeReconnection(ctx context.Context, url string) {
loop:
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-sys.Shutdown():
			return
		case <-ctx.Done():
			return
		case <-ws.shutdown:
			return
		case <-ws.reconnect:
			ws.logger.Warn("reconnecting...")
			ws.clearConnection()
			d, err := ws.dial()
			if err != nil {
				ws.logger.Errorf("ws connect to (%s), err: %+v", url, err)
				channel.TryPush(ws.reconnect, struct{}{})
				continue
			}

			ws.logger.Infof("ws connect to (%s) succeed, start subscribing...", url)
			ws.conn.Store(d)

			for _, register := range ws.registers {
				if err := ws.SendAndWait(ctx, register); err != nil {
					ws.clearConnection()
					ws.logger.Errorf("register, err: %+v", err)
					channel.TryPush(ws.reconnect, struct{}{})
					continue loop
				}
			}

			ws.logger.Info("ws subscribing succeed, connection available")

			go func() {
				defer d.Close()
				for {
					select {
					case <-sys.Shutdown():
						return
					case <-ctx.Done():
						return
					case <-ws.shutdown:
						return
					case <-d.Done():
						channel.TryPush(ws.reconnect, struct{}{})
						return
					case msg, ok := <-d.Message():
						if ok {
							ws.broadcast(msg)
						} else {
							channel.TryPush(ws.reconnect, struct{}{})
							return
						}
					}
				}
			}()
		}
	}
}

func (ws *WebSocket) broadcast(msg Message) {
	ws.subscribersLock.RLock()
	defer ws.subscribersLock.RUnlock()

	for _, ch := range ws.subscribers {
		if !channel.TryPush(ch, msg) {
			ws.logger.Warnf("broadcast message dropped for a subscriber, channel is full or closed")
		}
	}
}

func (ws *WebSocket) getConn() *dialing {
	d, ok := ws.conn.Load().(*dialing)
	if ok && d != nil {
		return d
	}

	return nil
}

func (ws *WebSocket) Start(ctx context.Context, registers ...Sidecar) {
	if ws.start.Swap(true) {
		return
	}

	ws.registers = registers

	go ws.observeReconnection(ctx, ws.url)

	channel.TryPush(ws.reconnect, struct{}{})
	timeout := time.After(DefaultTimeout)
	for {
		select {
		case <-sys.Shutdown():
			return
		case <-timeout:
			ws.logger.Warn("start ws timeout")
			return
		default:
			if ws.getConn() != nil {
				ws.logger.Info("start ws succeed")
				return
			}
		}
	}
}

func (ws *WebSocket) Close() {
	if ws.end.Swap(true) {
		return
	}

	channel.SafeClose(ws.shutdown)

	ws.subscribersLock.Lock()
	defer ws.subscribersLock.Unlock()
	for id, ch := range ws.subscribers {
		channel.SafeClose(ch)
		delete(ws.subscribers, id)
	}
}

func (ws *WebSocket) IsClose() bool {
	return channel.IsClose(ws.shutdown)
}

func (ws *WebSocket) Reconnect() {
	d, ok := ws.conn.Load().(*dialing)
	if ok && d != nil {
		d.Close()
	} else {
		channel.TryPush(ws.reconnect, struct{}{})
	}
}

func (ws *WebSocket) Subscribe() (<-chan Message, func()) {
	ws.subscribersLock.Lock()
	defer ws.subscribersLock.Unlock()

	id := ws.nextID.Add(1)
	ch := make(chan Message, _defaultMessageQueueCap)
	ws.subscribers[id] = ch

	unsubscribe := func() {
		ws.subscribersLock.Lock()
		defer ws.subscribersLock.Unlock()

		if ch, ok := ws.subscribers[id]; ok {
			channel.SafeClose(ch)
			delete(ws.subscribers, id)
		}
	}

	return ch, unsubscribe
}

func (ws *WebSocket) WriteJSON(v any, subscribeFunc ...bool) error {
	d := ws.getConn()
	if d == nil {
		return errors.Wrap(ErrConnectionClose, "nil ws connection")
	}

	return d.WriteJSON(v)
}

func (ws *WebSocket) WriteRaw(messageType MessageType, data []byte, subscribeFunc ...bool) error {
	d := ws.getConn()
	if d == nil {
		return errors.Wrap(ErrConnectionClose, "nil ws connection")
	}

	return d.WriteRaw(messageType, data)
}

func (ws *WebSocket) Len() int {
	ws.subscribersLock.RLock()
	defer ws.subscribersLock.RUnlock()

	return len(ws.subscribers)
}
