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
	message   chan Message
	shutdown  chan struct{}
	reconnect chan struct{}

	subscribe     []func(*dialing) error
	subscribeLock sync.RWMutex

	start  atomic.Bool
	end    atomic.Bool
	logger logs.Logger
}

func New(ctx context.Context, url string, ping ...bool) *WebSocket {
	ws := &WebSocket{
		url: url,
		dial: func() (*dialing, error) {
			return dial(ctx, url, ping...)
		},
		shutdown:  make(chan struct{}),
		reconnect: make(chan struct{}, 1),
		message:   make(chan Message, _defaultMessageQueueCap),
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

const DefaultWaitingMessageTimeout = 15 * time.Second

func SendAndWait[T any](ctx context.Context, ws *WebSocket, send func(context.Context, *WebSocket) error, isWaitTarget func(context.Context, T) bool, timeout ...time.Duration) (T, error) {
	done := make(chan error, 1)
	defer channel.SafeClose(done)

	waitTimeout := DefaultWaitingMessageTimeout
	if len(timeout) != 0 && timeout[0] > 0 {
		waitTimeout = timeout[0]
	}

	ctx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()

	var resp T
	go func() {
		for {
			select {
			case msg := <-ws.Message():
				resp, ok := ReadMessage[T](msg)
				if !ok {
					continue
				}

				if isWaitTarget(ctx, resp) {
					done <- nil
				}
			case <-ctx.Done():
				done <- ctx.Err()
				return
			case <-sys.Shutdown():
				done <- context.Canceled
				return
			}
		}
	}()

	if err := send(ctx, ws); err != nil {
		return resp, errors.Wrap(err, "send func error before waiting")
	}

	return resp, errors.Wrap(<-done, "websocket message")
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
			var d *dialing
			ws.conn.Store(d)
			d, err := ws.dial()
			if err != nil {
				ws.logger.Errorf("ws connect to (%s), err: %+v", url, err)
				channel.TryPush(ws.reconnect, struct{}{})
				continue
			}

			ws.logger.Info("ws connect to (%s) succeed, start subscribing...", url)

			ws.subscribeLock.RLock()
			subscribe := make([]func(*dialing) error, len(ws.subscribe))
			copy(subscribe, ws.subscribe)
			ws.subscribeLock.RUnlock()

			for _, fn := range subscribe {
				if err := fn(d); err != nil {
					ws.logger.Errorf("subscribe topic, err: %+v", err)
					channel.TryPush(ws.reconnect, struct{}{})
					continue loop
				}
			}

			ws.logger.Info("ws subscribing succeed, connection available")
			ws.conn.Store(d)

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
							ws.message <- msg
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

func (ws *WebSocket) getConn() *dialing {
	d, ok := ws.conn.Load().(*dialing)
	if ok && d != nil {
		return d
	}

	return nil
}

func (ws *WebSocket) Start(ctx context.Context) {
	if ws.start.Swap(true) {
		return
	}

	go ws.observeReconnection(ctx, ws.url)

	channel.TryPush(ws.reconnect, struct{}{})
	timeout := time.After(5 * time.Second)
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
}

func (ws *WebSocket) IsClose() bool {
	return channel.IsClose(ws.shutdown)
}

func (ws *WebSocket) ReconnectAndSubscribe() {
	d, ok := ws.conn.Load().(*dialing)
	if ok && d != nil {
		d.Close()
	} else {
		channel.TryPush(ws.reconnect, struct{}{})
	}
}

func (ws *WebSocket) Message() <-chan Message {
	return ws.message
}

func (ws *WebSocket) Produce() <-chan Message {
	return ws.message
}

func (ws *WebSocket) WriteJSON(v any, subscribeFunc ...bool) error {
	if len(subscribeFunc) != 0 && subscribeFunc[0] {
		ws.subscribeLock.Lock()
		ws.subscribe = append(ws.subscribe, func(d *dialing) error {
			return d.WriteJSON(v)
		})
		ws.subscribeLock.Unlock()
	}

	d := ws.getConn()
	if d == nil {
		return errors.Wrap(ErrConnectionClose, "nil ws connection")
	}

	return d.WriteJSON(v)
}

func (ws *WebSocket) WriteRaw(messageType MessageType, data []byte, subscribeFunc ...bool) error {
	if len(subscribeFunc) != 0 && subscribeFunc[0] {
		ws.subscribeLock.Lock()
		ws.subscribe = append(ws.subscribe, func(d *dialing) error {
			return d.WriteRaw(messageType, data)
		})
		ws.subscribeLock.Unlock()
	}

	d := ws.getConn()
	if d == nil {
		return errors.Wrap(ErrConnectionClose, "nil ws connection")
	}

	return d.WriteRaw(messageType, data)
}
