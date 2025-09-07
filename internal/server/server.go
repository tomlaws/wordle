package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/pkg/utils"
)

var Upgrader = websocket.Upgrader{}

func handleRead(client *Client) {
	defer func() {
		client.conn.Close()
	}()
	for {
		var msg json.RawMessage
		if err := client.conn.ReadJSON(&msg); err != nil {
			client.error <- err
			log.Printf("Error reading message from player %s: %v", client.nickname, err)
			break
		}
		client.incoming <- msg
		log.Printf("Received message from player %s: %+v", client.nickname, utils.JsonToString(msg))
	}
}

func handleWrite(client *Client) {
	defer func() {
		client.conn.Close()
	}()
	for {
		msg := <-client.outgoing
		if err := client.conn.WriteJSON(msg); err != nil {
			client.error <- err
			log.Printf("Error sending message to player %s: %v", client.nickname, err)
			break
		}
		log.Printf("Sending message to player %s: %s", client.nickname, utils.JsonToString(msg))
	}
}

func socketHandler(newClientCallback func(client *Client)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.URL.Query().Get("nickname")
		// Reject if username is blank
		if nickname == "" {
			log.Printf("Player connected with blank nickname")
			return
		}

		// Upgrade our raw HTTP connection to a websocket based one
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}

		client := &Client{
			id:       uuid.New().String(),
			nickname: nickname,
			conn:     conn,
			incoming: make(chan json.RawMessage),
			outgoing: make(chan json.RawMessage),
			error:    make(chan error),
		}
		go handleRead(client)
		go handleWrite(client)
		newClientCallback(client)
	}
}

func NewServer(
	newClientCallback func(client *Client),
) func(w http.ResponseWriter, r *http.Request) {
	return socketHandler(newClientCallback)
}
