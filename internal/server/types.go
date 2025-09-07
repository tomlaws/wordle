package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	id       string
	nickname string
	incoming chan json.RawMessage
	outgoing chan json.RawMessage
	error    chan error
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Nickname() string {
	return c.nickname
}

func (c *Client) Incoming() chan json.RawMessage {
	return c.incoming
}

func (c *Client) Outgoing() chan json.RawMessage {
	return c.outgoing
}

func (c *Client) Err() chan error {
	return c.error
}
