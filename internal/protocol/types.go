package protocol

import "encoding/json"

type Protocol struct {
	registry map[MessageType]func() Payload
}

type MessageType string

type Message struct {
	Type    MessageType
	Payload json.RawMessage
}

type Payload interface {
	MessageType() MessageType
}
