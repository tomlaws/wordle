package server

import (
	"encoding/json"

	"github.com/tomlaws/wordle/internal/game"
)

const (
	// Server to Client
	MsgTypeGameStart   = "game_start"
	MsgTypeInvalidWord = "invalid_word"
	MsgTypeFeedback    = "feedback"
	MsgTypeGameOver    = "game_over"
	// Client to Server
	MsgTypeGuess       = "guess"
	MsgTypeConfirmPlay = "confirm_play"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type GameStartPayload struct {
	MaxAttempts int `json:"max_attempts"`
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
	Won    bool   `json:"won"`
	Answer string `json:"answer"`
}

type ConfirmPlayPayload struct {
	Confirm bool `json:"confirm"`
}
