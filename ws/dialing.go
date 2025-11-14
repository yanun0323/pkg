package ws

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yanun0323/errors"
	"github.com/yanun0323/pkg/channel"
	"github.com/yanun0323/pkg/sys"
)

const (
	_symbolRecv        string = "\x1b[32m⬇\x1b[0m"
	_symbolWrite       string = "\x1b[33m⬆\x1b[0m"
	_debugMessageLimit int    = 100
)

var (
	Debug              = false
	ErrNilInstance     = errors.New("nil instance")
	ErrConnectionClose = errors.New("connection closed")
)

type dialing struct {
	conn    *websocket.Conn
	message chan Message
	writeMu sync.Mutex
	done    chan struct{}
	close   chan struct{}
	logger  *slog.Logger
}

// dial creates
func dial(ctx context.Context, url string, ping ...bool) (*dialing, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "dial to (%s)", url)
	}

	if conn == nil {
		return nil, errors.Wrapf(ErrNilInstance, "dial to (%s)", url)
	}

	pong := make(chan struct{}, 1)
	d := &dialing{
		conn:    conn,
		message: make(chan Message, _defaultMessageQueueCap),
		writeMu: sync.Mutex{},
		done:    make(chan struct{}, 1),
		close:   make(chan struct{}),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})).With(
			slog.String("dialing", url),
			slog.String("id", newLogID()),
		),
	}

	// receiver
	go func() {
		defer channel.SafeClose(d.done)
		defer channel.SafeClose(d.message)
		for {
			select {
			case <-sys.Shutdown():
				return
			default:
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					d.logger.Error(fmt.Sprintf("read message, err: %+v", err))
					d.logger.Error("stop reading message")

					return
				}

				mt := MessageType(messageType)
				switch mt {
				case MessageTypeClose:
					if Debug {
						d.logger.Info(fmt.Sprintf("%s close message: %s", _symbolRecv, string(message)))
					}
					return
				case MessageTypePong:
					channel.TryPush(pong, struct{}{})
					if Debug {
						d.logger.Info(fmt.Sprintf("%s pong message", _symbolRecv))
					}
				default:
					if Debug {
						msg := fmt.Sprintf("%s %s message: %s", _symbolRecv, mt, string(message))
						if len(msg) >= _debugMessageLimit {
							msg = msg[:_debugMessageLimit] + "..."
						}
						d.logger.Info(msg)
					}
				}

				d.message <- Message{
					Type: mt,
					Data: message,
				}
			}
		}
	}()

	// pinging & checking pong
	if len(ping) != 0 && ping[0] {
		go func() {
			defer channel.SafeClose(d.done)
			defer channel.SafeClose(pong)
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			channel.TryPush(pong, struct{}{})

			for {
				select {
				case <-sys.Shutdown():
					return
				case <-d.done:
					return
				case <-ticker.C:
					if _, ok := channel.TryReceive(pong); !ok {
						d.logger.Error("receive no pong message")
						return
					}

					d.writeMu.Lock()
					err := conn.WriteMessage(MessageTypePing.Int(), nil)
					d.writeMu.Unlock()
					if err != nil {
						d.logger.Error(fmt.Sprintf("ping, err: %+v", err))
						return
					}

					d.logger.Debug("ping succeed")
				}
			}
		}()
	}

	// closing connection
	go func() {
		defer channel.SafeClose(d.done)
		defer channel.SafeClose(pong)
		select {
		case <-sys.Shutdown():
		case <-d.done:
		case <-ctx.Done():
		case <-d.close:
		}

		if err := conn.Close(); err != nil {
			d.logger.Error(fmt.Sprintf("closing dialing, err: %+v", err))
		} else {
			d.logger.Info("dialing closed")
		}
	}()

	return d, nil
}

func (c *dialing) IsClose() bool {
	return channel.IsClose(c.done)
}

func (c *dialing) Done() <-chan struct{} {
	return c.done
}

func (c *dialing) Close() {
	if !c.IsClose() {
		channel.SafeClose(c.close)
	}
}

func (c *dialing) Message() <-chan Message {
	return c.message
}

func (c *dialing) WriteJSON(v any) error {
	if c.IsClose() {
		return errors.Wrap(ErrConnectionClose, "dialing closed")
	}

	c.writeMu.Lock()
	err := c.conn.WriteJSON(v)
	c.writeMu.Unlock()
	if err != nil {
		return errors.Wrapf(err, "write json (%v)", v)
	}

	if Debug {
		msg := fmt.Sprintf("%s write json (%v)", _symbolWrite, v)
		if len(msg) >= _debugMessageLimit {
			msg = msg[:_debugMessageLimit] + "..."
		}

		c.logger.Info(msg)
	}

	return nil
}

func (c *dialing) WriteRaw(messageType MessageType, data []byte) error {
	if c.IsClose() {
		return errors.Wrap(ErrConnectionClose, "dialing closed")
	}

	c.writeMu.Lock()
	err := c.conn.WriteMessage(messageType.Int(), data)
	c.writeMu.Unlock()
	if err != nil {
		return errors.Wrapf(err, "write message. type(%s) data(%s)",
			messageType,
			string(data),
		)
	}

	if Debug {
		msg := fmt.Sprintf("%s write raw [%s](%s)", _symbolWrite, messageType, string(data))
		if len(msg) >= _debugMessageLimit {
			msg = msg[:_debugMessageLimit] + "..."
		}

		c.logger.Info(msg)
	}

	return nil
}
