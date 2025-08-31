package client

import (
	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/server"
)

type Client struct {
	conn         *websocket.Conn
	input        chan Input
	inputTrigger chan InputTrigger
	incoming     chan *server.Message
	outgoing     chan *server.Message
	error        chan error
}

type InputCategory int

const (
	GuessWord InputCategory = iota
	PlayAgain
)

type InputTrigger struct {
	Category InputCategory
}

type Input struct {
	Category InputCategory
	Text     string
}
