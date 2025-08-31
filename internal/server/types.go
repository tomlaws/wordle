package server

import (
	"encoding/json"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
)

const (
	// Server to Client
	MsgTypePlayerInfo       = "player_info"
	MsgTypeMatching         = "matching"
	MsgTypeGameStart        = "game_start"
	MsgTypeRoundStart       = "round_start"
	MsgTypeInvalidWord      = "invalid_word"
	MsgTypeFeedback         = "feedback"
	MsgTypeGameOver         = "game_over"
	MsgTypePlayAgainTimeout = "play_again_timeout"
	MsgTypeGuessTimeout     = "guess_timeout"
	// Client to Server
	MsgTypeTyping    = "typing"
	MsgTypeGuess     = "guess"
	MsgTypePlayAgain = "play_again"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type GameStartPayload struct {
	MaxGuesses int     `json:"max_guesses"`
	Player1    *Player `json:"player1"`
	Player2    *Player `json:"player2"`
}

type RoundStartPayload struct {
	Player  *Player `json:"player"`
	Round   int     `json:"round"`
	Timeout int     `json:"timeout"` // seconds
}

type InvalidWordPayload struct {
	Player *Player `json:"player"`
	Word   string  `json:"word"`
}

type TypingPayload struct {
	Word string `json:"word"`
}

type GuessPayload struct {
	Word string `json:"word"`
}

type FeedbackPayload struct {
	Player   *Player             `json:"player"`
	Round    int                 `json:"round"`
	Feedback []game.LetterResult `json:"feedback"`
}

type GameOverPayload struct {
	Winner *Player `json:"winner"`
	Answer string  `json:"answer"`
}

type GuessTimeoutPayload struct {
	Player *Player `json:"player"`
}

type PlayAgainPayload struct {
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
