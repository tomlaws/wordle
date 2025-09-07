package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func NewClient(ipAddress string, nickname string) (*Client, error) {
	url := url.URL{Scheme: "ws", Host: ipAddress, Path: "/socket", RawQuery: fmt.Sprintf("nickname=%s", nickname)}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	client := &Client{
		url:      url,
		conn:     conn,
		incoming: make(chan json.RawMessage),
		outgoing: make(chan json.RawMessage),
		error:    make(chan error),
	}
	go handleRead(client)
	go handleWrite(client)
	return client, nil
}

func handleRead(client *Client) {
	defer func() {
		close(client.incoming)
	}()
	for {
		var msg json.RawMessage
		if err := client.conn.ReadJSON(&msg); err != nil {
			client.error <- err
			// log.Printf("Error reading message from server: %v", err)
			break
		}
		client.incoming <- msg
		// log.Printf("Received message from server: %+v", utils.JsonToString(msg))
	}
}

func handleWrite(client *Client) {
	defer func() {
		close(client.outgoing)
	}()
	for {
		msg := <-client.outgoing
		if err := client.conn.WriteJSON(msg); err != nil {
			client.error <- err
			// log.Printf("Error sending message to server: %v", err)
			break
		}
		// log.Printf("Sending message to server: %s", utils.JsonToString(msg))
	}
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

func (c *Client) Stop() {
	c.conn.Close()
}
