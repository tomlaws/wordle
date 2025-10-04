package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/pkg/utils"
)

var Upgrader = websocket.Upgrader{}

const pingInterval = 15 * time.Second
const pongWait = 25 * time.Second
const writeWait = 5 * time.Second

func handleRead(client *Client) {
	defer func() {
		client.conn.Close()
	}()
	// handling pong messages from client
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		// log.Printf("Received pong from player %s", client.nickname)
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
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
	ticker := time.NewTicker(pingInterval)
	defer func() {
		client.conn.Close()
		ticker.Stop()
	}()
	for {
		select {
		case msg := <-client.outgoing:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteJSON(msg); err != nil {
				client.error <- err
				log.Printf("Error sending message to player %s: %v", client.nickname, err)
				break
			}
			log.Printf("Sending message to player %s: %s", client.nickname, utils.JsonToString(msg))
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.error <- err
				log.Printf("Error sending ping to player %s: %v", client.nickname, err)
			}
		}
	}
}

func socketHandler(newClientCallback func(client *Client)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := strings.TrimSpace(r.URL.Query().Get("nickname"))
		if len(nickname) < 3 || len(nickname) > 16 {
			log.Printf("Player connected with invalid nickname length: %s", nickname)
			http.Error(w, "Nickname must be between 3 and 16 characters", http.StatusBadRequest)
			return
		}
		Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
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
	// set deadline for both read and write to 60 seconds
	return socketHandler(newClientCallback)
}
