package protocol

import "encoding/json"

type Protocol struct {
	registry map[MessageType]func() Payload
}

type MessageType string

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Payload interface {
	MessageType() MessageType
}
