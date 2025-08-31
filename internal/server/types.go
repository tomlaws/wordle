package server

import (
	"encoding/json"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
)

const (
	// Server to Client
	MsgTypePlayerInfo  = "player_info"
	MsgTypeMatching    = "matching"
	MsgTypeGameStart   = "game_start"
	MsgTypeTurnStart   = "turn_start"
	MsgTypeInvalidWord = "invalid_word"
	MsgTypeFeedback    = "feedback"
	MsgTypeGameOver    = "game_over"
	// Client to Server
	MsgTypeTyping    = "typing"
	MsgTypeGuess     = "guess"
	MsgTypePlayAgain = "play_again"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type GameStartPayload struct {
	MaxAttempts int     `json:"max_attempts"`
	Player1     *Player `json:"player1"`
	Player2     *Player `json:"player2"`
}

type TurnStartPayload struct {
	Player *Player `json:"player"`
}

type TypingPayload struct {
	Word string `json:"word"`
}

type GuessRequest struct {
	Word string `json:"word"`
}

type InvalidWordResponse struct {
	Word string `json:"word"`
}

type FeedbackResponse struct {
	Feedback []game.LetterResult `json:"feedback"`
	Round    int                 `json:"round"`
}

type GameOverPayload struct {
	Winner *Player `json:"winner"`
	Answer string  `json:"answer"`
}

type ConfirmPlayPayload struct {
	Confirm bool `json:"confirm"`
}

type PlayerState int

const (
	Disconnected PlayerState = iota
	Connected
	InGame
)

type Player struct {
	conn     *websocket.Conn
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	state    atomic.Int32
	incoming chan *Message
	outgoing chan *Message
	error    chan error
}
