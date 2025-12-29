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
	_defaultBackoffMin          = 100 * time.Millisecond
	_defaultBackoffMax          = 30 * time.Second
	_defaultBackoffFactor       = 2.0
)

// Option defines websocket settings for New.
type Option struct {
	// Ping enables an automatic ping goroutine to keep the connection alive.
	Ping bool
	// BackoffOption defines reconnection backoff behavior.
	Backoff BackoffOption
}

// BackoffOption defines reconnection backoff behavior.
type BackoffOption struct {
	// Min is the minimum delay before a reconnect attempt.
	Min time.Duration
	// Max is the maximum delay before a reconnect attempt.
	Max time.Duration
	// Factor is the multiplier for exponential backoff growth.
	Factor float64
}

type backoffState struct {
	min     time.Duration
	max     time.Duration
	factor  float64
	current time.Duration
}

// WebSocket holds a websocket connection and treat it like a producer.
//
// It handles all clients as consumer, publish message data from websocket to all consumers.
//
// It also handles reconnection automatically.
type WebSocket struct {
	conn atomic.Value // *dialing
	url  string

	dial      func() (*dialing, error)
	shutdown  chan struct{}
	reconnect chan struct{}

	registersLock sync.RWMutex
	registers     []Sidecar

	subscribers     map[uint64]chan Message
	subscribersLock sync.RWMutex
	nextID          atomic.Uint64

	start  atomic.Bool
	end    atomic.Bool
	option Option
	logger logs.Logger
}

// Sidecar defines hooks and timing for SendAndWait execution.
type Sidecar struct {
	// Sender emits a request through the websocket.
	Sender func(context.Context, *WebSocket) error
	// Waiter inspects incoming messages and returns true when the expected response arrives.
	Waiter func(context.Context, Message) (isExpected bool, failure error)
	// Timeout limits how long SendAndWait waits for an expected response.
	Timeout time.Duration
}

// New creates a new websocket connection without connecting to the websocket.
func New(ctx context.Context, url string, opts ...Option) *WebSocket {
	option := Option{}
	if len(opts) > 0 {
		option = opts[0]
	}
	option = normalizeOption(option)

	ws := &WebSocket{
		url: url,
		dial: func() (*dialing, error) {
			return dial(ctx, url, option.Ping)
		},
		shutdown:    make(chan struct{}),
		reconnect:   make(chan struct{}, 1),
		subscribers: make(map[uint64]chan Message),
		option:      option,
		logger: logs.Get(ctx).With(
			"websocket", newLogID(),
			"url", url,
		),
	}

	return ws
}

// ReadMessage parses the message with provided types
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
	// DefaultStartTimeout is the default timeout for WebSocket.Start(...)
	DefaultStartTimeout = 15 * time.Second
	// DefaultWaitTimeout is the default timeout for WebSocket.SendAndWait(...)
	DefaultWaitTimeout = 15 * time.Second
)

// SendAndWait using a Sidecar to send a message and wait for particular response
func (ws *WebSocket) SendAndWait(ctx context.Context, executor Sidecar, appendIntoRegister ...bool) error {
	if executor.Sender == nil || executor.Waiter == nil {
		return errors.New("invalid hook, require sender and waiter")
	}

	if len(appendIntoRegister) != 0 && appendIntoRegister[0] {
		ws.registersLock.Lock()
		ws.registers = append(ws.registers, executor)
		ws.registersLock.Unlock()
	}

	done := make(chan error, 1)
	defer channel.SafeClose(done)

	waitTimeout := executor.Timeout
	if executor.Timeout <= 0 {
		waitTimeout = DefaultWaitTimeout
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

	if err := <-done; err != nil {
		return errors.Wrap(err, "unexpected result")
	}

	return nil
}

func (ws *WebSocket) clearConnection() {
	var d *dialing
	ws.conn.Store(d)
}

func normalizeOption(option Option) Option {
	option.Backoff = normalizeBackoff(option.Backoff)
	return option
}

func normalizeBackoff(option BackoffOption) BackoffOption {
	if option.Min <= 0 {
		option.Min = _defaultBackoffMin
	}
	if option.Max <= 0 {
		option.Max = _defaultBackoffMax
	}
	if option.Max < option.Min {
		option.Max = option.Min
	}
	if option.Factor <= 1 {
		option.Factor = _defaultBackoffFactor
	}

	return option
}

func newBackoffState(option BackoffOption) backoffState {
	option = normalizeBackoff(option)
	return backoffState{
		min:    option.Min,
		max:    option.Max,
		factor: option.Factor,
	}
}

func (b *backoffState) Next() time.Duration {
	if b.min <= 0 {
		return 0
	}
	if b.current <= 0 {
		b.current = b.min
		return b.current
	}

	next := time.Duration(float64(b.current) * b.factor)
	if next < b.min {
		next = b.min
	}
	if b.max > 0 && next > b.max {
		next = b.max
	}
	b.current = next
	return b.current
}

func (b *backoffState) Reset() {
	b.current = 0
}

func (ws *WebSocket) waitBackoff(ctx context.Context, delay time.Duration) bool {
	if delay <= 0 {
		return true
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-sys.Shutdown():
		return false
	case <-ctx.Done():
		return false
	case <-ws.shutdown:
		return false
	case <-timer.C:
		return true
	}
}

func (ws *WebSocket) observeReconnection(ctx context.Context, url string, start chan struct{}) {
	var once sync.Once
	backoff := newBackoffState(ws.option.Backoff)

loop:
	for {
		select {
		case <-sys.Shutdown():
			return
		case <-ctx.Done():
			return
		case <-ws.shutdown:
			return
		case <-ws.reconnect:
			dur := backoff.Next()
			ws.logger.Warnf("reconnecting in %s...", dur)
			if !ws.waitBackoff(ctx, dur) {
				return
			}

			ws.logger.Warn("reconnecting...")
			ws.clearConnection()
			d, err := ws.dial()
			if err != nil {
				ws.logger.Errorf("ws connect to (%s), err: %+v", url, err)
				channel.TryPush(ws.reconnect, struct{}{})
				continue
			}

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

			ws.logger.Infof("ws connect to (%s) succeed, start subscribing...", url)
			ws.conn.Store(d)

			ws.registersLock.RLock()
			for _, register := range ws.registers {
				if err := ws.SendAndWait(ctx, register); err != nil {
					ws.clearConnection()
					ws.logger.Errorf("register, err: %+v", err)
					channel.TryPush(ws.reconnect, struct{}{})
					continue loop
				}
			}
			ws.registersLock.RUnlock()

			once.Do(func() {
				channel.SafeClose(start)
			})

			backoff.Reset()
			ws.logger.Info("ws subscribing succeed, connection available")
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

// Start starts connecting to the websocket.
//
// Args:
//   - register: represents operations which must be invoked after every websocket connecting/reconnecting
func (ws *WebSocket) Start(_ctx context.Context, registers ...Sidecar) error {
	if ws.start.Swap(true) {
		return nil
	}

	ws.registers = registers

	start := make(chan struct{})
	defer channel.SafeClose(start)

	go ws.observeReconnection(_ctx, ws.url, start)

	channel.TryPush(ws.reconnect, struct{}{})
	ctx, cancel := context.WithTimeout(_ctx, DefaultStartTimeout)
	defer cancel()

	select {
	case <-sys.Shutdown():
		return errors.Wrap(context.Canceled)
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "timeout")
	case <-start:
		ws.logger.Info("start ws succeed")
	}

	return nil
}

// Close closes the websocket connection
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

// IsClose returns whether the websocket is closed or not
func (ws *WebSocket) IsClose() bool {
	return channel.IsClose(ws.shutdown)
}

// Reconnect tries to reconnect the websocket
func (ws *WebSocket) Reconnect() {
	d, ok := ws.conn.Load().(*dialing)
	if ok && d != nil {
		d.Close()
	} else {
		channel.TryPush(ws.reconnect, struct{}{})
	}
}

// Subscribe subscribes the websocket message producer
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

// WriteJSON writes a JSON message to the websocket connection
func (ws *WebSocket) WriteJSON(v any) error {
	d := ws.getConn()
	if d == nil {
		return errors.Wrap(ErrConnectionClose, "nil ws connection")
	}

	return d.WriteJSON(v)
}

// WriteRaw writes a raw message to the websocket connection
func (ws *WebSocket) WriteRaw(messageType MessageType, data []byte, subscribeFunc ...bool) error {
	d := ws.getConn()
	if d == nil {
		return errors.Wrap(ErrConnectionClose, "nil ws connection")
	}

	return d.WriteRaw(messageType, data)
}

// Len returns the subscribers number of websocket connection
func (ws *WebSocket) Len() int {
	ws.subscribersLock.RLock()
	defer ws.subscribersLock.RUnlock()

	return len(ws.subscribers)
}
