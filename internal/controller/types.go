package controller

import "github.com/tomlaws/wordle/internal/protocol"

type Controller struct {
	input        chan Input
	inputTrigger chan InputTrigger
	incoming     <-chan protocol.Payload
	outgoing     chan<- protocol.Payload
	error        chan error
}

type InputTrigger struct {
	Category InputCategory
}

type Input struct {
	Category InputCategory
	Text     string
}

type InputCategory int

const (
	GuessWord InputCategory = iota
	PlayAgain
)

type Me struct {
	ID       string
	Nickname string
}
