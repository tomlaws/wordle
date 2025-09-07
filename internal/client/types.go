package client

import (
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	url      url.URL
	conn     *websocket.Conn
	incoming chan json.RawMessage
	outgoing chan json.RawMessage
	error    chan error
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Payload interface {
	MessageType() string
}
