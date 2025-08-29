package server

import (
	"encoding/json"

	"github.com/tomlaws/wordle/internal/game"
)

const (
	MsgTypeGameStart   = "game_start"
	MsgTypeGameOver    = "game_over"
	MsgTypeGuess       = "guess"
	MsgTypeInvalidWord = "invalid_word"
	MsgTypeFeedback    = "feedback"
	MsgTypeConfirmPlay = "confirm_play"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type GameStartPayload struct {
	MaxAttempts int `json:"max_attempts"`
}

type GameOverPayload struct {
	Won    bool   `json:"won"`
	Answer string `json:"answer"`
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

type ConfirmPlayPayload struct {
	Confirm bool `json:"confirm"`
}
