package ws

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/yanun0323/errors"
)

type MessageType int

const (
	MessageTypeText   MessageType = websocket.TextMessage
	MessageTypeBinary MessageType = websocket.BinaryMessage
	MessageTypeClose  MessageType = websocket.CloseMessage
	MessageTypePing   MessageType = websocket.PingMessage
	MessageTypePong   MessageType = websocket.PongMessage
)

func (mt MessageType) Int() int {
	return int(mt)
}

func (mt MessageType) String() string {
	switch mt {
	case MessageTypeText:
		return "Text"
	case MessageTypeBinary:
		return "Binary"
	case MessageTypeClose:
		return "Close"
	case MessageTypePing:
		return "Ping"
	case MessageTypePong:
		return "Pong"
	default:
		return fmt.Sprintf("Unknown(%d)", mt)
	}
}

type Message struct {
	Type MessageType
	Data []byte
}

func (m Message) String() string {
	return fmt.Sprintf("type: %s, data: %s", m.Type, m.Data)
}

func (m Message) Unmarshal(p any) error {
	return errors.Wrap(json.Unmarshal(m.Data, p), "unmarshal data")
}
