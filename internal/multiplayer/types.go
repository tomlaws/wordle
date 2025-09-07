package multiplayer

import (
	"encoding/json"

	"github.com/tomlaws/wordle/internal/game"
	"github.com/tomlaws/wordle/internal/protocol"
)

type Client interface {
	ID() string
	Nickname() string
	Incoming() chan json.RawMessage
	Outgoing() chan json.RawMessage
	Err() chan error
}

type Player struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	incoming chan protocol.Payload
	outgoing chan protocol.Payload
	err      chan error
}

type Lobby struct {
	wordList   *game.WordList
	maxGuesses int
	thinkTime  int
	queue      chan *Player
}

const (
	MsgTypeTyping    protocol.MessageType = "typing"
	MsgTypeGuess     protocol.MessageType = "guess"
	MsgTypePlayAgain protocol.MessageType = "play_again"
)

const (
	MsgTypePlayerInfo   protocol.MessageType = "player_info"
	MsgTypeMatching     protocol.MessageType = "matching"
	MsgTypeGameStart    protocol.MessageType = "game_start"
	MsgTypeRoundStart   protocol.MessageType = "round_start"
	MsgTypeInvalidWord  protocol.MessageType = "invalid_word"
	MsgTypeGuessTimeout protocol.MessageType = "guess_timeout"
	MsgTypeFeedback     protocol.MessageType = "feedback"
	MsgTypeGameOver     protocol.MessageType = "game_over"
)

var PayloadRegistry = map[protocol.MessageType]func() protocol.Payload{
	MsgTypePlayerInfo:   func() protocol.Payload { return &PlayerInfoPayload{} },
	MsgTypeMatching:     func() protocol.Payload { return &MatchingPayload{} },
	MsgTypeGameStart:    func() protocol.Payload { return &GameStartPayload{} },
	MsgTypeRoundStart:   func() protocol.Payload { return &RoundStartPayload{} },
	MsgTypeInvalidWord:  func() protocol.Payload { return &InvalidWordPayload{} },
	MsgTypeGuessTimeout: func() protocol.Payload { return &GuessTimeoutPayload{} },
	MsgTypeFeedback:     func() protocol.Payload { return &FeedbackPayload{} },
	MsgTypeGameOver:     func() protocol.Payload { return &GameOverPayload{} },

	MsgTypeTyping:    func() protocol.Payload { return &TypingPayload{} },
	MsgTypeGuess:     func() protocol.Payload { return &GuessPayload{} },
	MsgTypePlayAgain: func() protocol.Payload { return &PlayAgainPayload{} },
}

type PlayerInfoPayload struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

func (p *PlayerInfoPayload) MessageType() protocol.MessageType {
	return MsgTypePlayerInfo
}

type MatchingPayload struct {
}

func (p *MatchingPayload) MessageType() protocol.MessageType {
	return MsgTypeMatching
}

type GameStartPayload struct {
	MaxGuesses int     `json:"max_guesses"`
	Player1    *Player `json:"player1"`
	Player2    *Player `json:"player2"`
}

func (p *GameStartPayload) MessageType() protocol.MessageType {
	return MsgTypeGameStart
}

type RoundStartPayload struct {
	Player  *Player `json:"player"`
	Round   int     `json:"round"`
	Timeout int     `json:"timeout"` // seconds
}

func (p *RoundStartPayload) MessageType() protocol.MessageType {
	return MsgTypeRoundStart
}

type InvalidWordPayload struct {
	Player *Player `json:"player"`
	Word   string  `json:"word"`
}

func (p *InvalidWordPayload) MessageType() protocol.MessageType {
	return MsgTypeInvalidWord
}

type GuessTimeoutPayload struct {
	Player *Player `json:"player"`
}

func (p *GuessTimeoutPayload) MessageType() protocol.MessageType {
	return MsgTypeGuessTimeout
}

type FeedbackPayload struct {
	Player   *Player             `json:"player"`
	Round    int                 `json:"round"`
	Feedback []game.LetterResult `json:"feedback"`
}

func (p *FeedbackPayload) MessageType() protocol.MessageType {
	return MsgTypeFeedback
}

type GameOverPayload struct {
	Winner *Player `json:"winner"`
	Answer string  `json:"answer"`
}

func (p *GameOverPayload) MessageType() protocol.MessageType {
	return MsgTypeGameOver
}

type TypingPayload struct {
	Player *Player `json:"player"`
	Word   string  `json:"word"`
}

func (p *TypingPayload) MessageType() protocol.MessageType {
	return MsgTypeTyping
}

type GuessPayload struct {
	Word string `json:"word"`
}

func (p *GuessPayload) MessageType() protocol.MessageType {
	return MsgTypeGuess
}

type PlayAgainPayload struct {
	Confirm bool `json:"confirm"`
}

func (p *PlayAgainPayload) MessageType() protocol.MessageType {
	return MsgTypePlayAgain
}
