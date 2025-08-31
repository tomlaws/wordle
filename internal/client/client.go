package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
	"github.com/tomlaws/wordle/internal/server"
)

func New(ipAddress string, nickname string) (*Client, error) {
	// Connect to the websocket (port 8080)
	u := url.URL{
		Scheme:   "ws",
		Host:     fmt.Sprintf("%s:%d", ipAddress, 8080),
		Path:     "/socket",
		RawQuery: fmt.Sprintf("nickname=%s", nickname),
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:         conn,
		input:        make(chan Input),
		inputTrigger: make(chan InputTrigger),
		incoming:     make(chan *server.Message),
		outgoing:     make(chan *server.Message),
		error:        make(chan error),
	}, nil
}

func handleInputUnderlying(lines chan interface{}) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines <- s.Text()
	}
	lines <- s.Err()
}

func handleInput(client *Client) {
	input := make(chan interface{})
	go handleInputUnderlying(input)
	var tr *InputTrigger = nil
	for {
		select {
		case trigger := <-client.inputTrigger:
			tr = &trigger
		case line := <-input:
			if tr != nil {
				client.input <- Input{Category: tr.Category, Text: line.(string)}
				tr = nil
			}
		}
	}
}

func handleRead(client *Client) {
	defer client.Stop()
	for {
		var msg server.Message
		if err := client.conn.ReadJSON(&msg); err != nil {
			client.error <- err
			// log.Printf("Error reading message %v", err)
			break
		}
		client.incoming <- &msg
	}
}

func handleWrite(client *Client) {
	defer client.Stop()
	for {
		msg := <-client.outgoing
		// log.Printf("Sending message to server: %s", msg.Type)
		if err := client.conn.WriteJSON(msg); err != nil {
			client.error <- err
			// log.Printf("Error sending message %v", err)
			break
		}
	}
}

func (c *Client) Stop() {
	c.conn.Close()
	os.Exit(0)
}

func (c *Client) Start(input io.Reader, output io.Writer) error {
	defer c.Stop()
	go handleInput(c)
	go handleRead(c)
	go handleWrite(c)
	var err error
	var me server.Player
	var maxGuesses int
	var currentRound int
	var isOddPlayer bool
	for {
		select {
		case msg := <-c.incoming:
			switch msg.Type {
			case server.MsgTypePlayerInfo:
				if err := json.Unmarshal(msg.Payload, &me); err != nil {
					log.Println("Error during player info payload unmarshalling:", err)
					return err
				}
				fmt.Fprintf(output, "Welcome to Wordle, %s!\n", me.Nickname)
			case server.MsgTypeMatching:
				fmt.Fprintf(output, "Finding opponent...\n")
			case server.MsgTypeGameStart:
				// Handle game start
				var gameStartPayload server.GameStartPayload
				if err := json.Unmarshal(msg.Payload, &gameStartPayload); err != nil {
					log.Println("Error during game start payload unmarshalling:", err)
					return err
				}
				currentRound = 1
				maxGuesses = gameStartPayload.MaxGuesses
				isOddPlayer = gameStartPayload.Player1.ID == me.ID
				var opponent *server.Player
				if isOddPlayer {
					opponent = gameStartPayload.Player2
				} else {
					opponent = gameStartPayload.Player1
				}
				fmt.Fprintf(output, "You are playing against %s\n", opponent.Nickname)
				fmt.Fprintln(output, "Guess the 5-letter word in", maxGuesses, "rounds.")
			case server.MsgTypeRoundStart:
				var roundStartPayload server.RoundStartPayload
				if err := json.Unmarshal(msg.Payload, &roundStartPayload); err != nil {
					log.Println("Error during round start payload unmarshalling:", err)
					return err
				}
				sameRound := roundStartPayload.Round == currentRound
				currentRound = roundStartPayload.Round
				timeout := roundStartPayload.Timeout
				// Handle guess input when it's the player's turn
				if roundStartPayload.Player.ID == me.ID {
					fmt.Fprintf(output, "=====Round (%d/%d)=====\n", currentRound, maxGuesses)
					c.inputTrigger <- InputTrigger{Category: GuessWord}
					if sameRound {
						fmt.Fprintln(output, "You have", timeout, "seconds to make your guess.")
					}
					fmt.Fprintf(output, "Enter your guess (%d/%d): ", currentRound, maxGuesses)
				} else {
					// Wait for opponent's guess
					fmt.Fprintf(output, "=====Round (%d/%d)=====\n", currentRound, maxGuesses)
					fmt.Fprintln(output, "Waiting for opponent's guess...")
				}
			case server.MsgTypeInvalidWord:
				//fmt.Fprintln(output, "Invalid word. Please try again.")
				var invalidWordPayload server.InvalidWordPayload
				if err := json.Unmarshal(msg.Payload, &invalidWordPayload); err != nil {
					log.Println("Error during invalid word payload unmarshalling:", err)
					return err
				}
				if invalidWordPayload.Player.ID == me.ID {
					fmt.Fprintln(output, "Invalid word. Please try again.")
				} else {
					fmt.Fprintf(output, "Opponent guessed an invalid word: %s\n", invalidWordPayload.Word)
				}
			case server.MsgTypeFeedback:
				var feedbackPayload server.FeedbackPayload
				if err := json.Unmarshal(msg.Payload, &feedbackPayload); err != nil {
					log.Println("Error during feedback payload unmarshalling:", err)
					return err
				}
				player := feedbackPayload.Player
				if player.ID == me.ID {
					fmt.Printf("You guessed: ")
				} else {
					fmt.Printf("Opponent guessed: ")
				}
				currentRound = feedbackPayload.Round + 1
				// Display feedback to the user
				for _, lr := range feedbackPayload.Feedback {
					switch lr.MatchType {
					case game.Hit:
						fmt.Fprintf(output, "[%c] ", lr.Letter)
					case game.Present:
						fmt.Fprintf(output, "(%c) ", lr.Letter)
					case game.Miss:
						fmt.Fprintf(output, " %c  ", lr.Letter)
					}
				}
				fmt.Fprintln(output)
			case server.MsgTypeGameOver:
				var gameOverPayload server.GameOverPayload
				if err := json.Unmarshal(msg.Payload, &gameOverPayload); err != nil {
					log.Println("Error during game over payload unmarshalling:", err)
					return err
				}
				if gameOverPayload.Winner == nil {
					fmt.Fprintln(output, "It's a draw! The correct word was:", gameOverPayload.Answer)
				} else if gameOverPayload.Winner.ID == me.ID {
					fmt.Fprintln(output, "Congratulations! You've won!")
				} else {
					fmt.Fprintln(output, "You've lost! The correct word was:", gameOverPayload.Answer)
				}

				// Ask for a new game
				fmt.Fprint(output, "Do you want to play again? (y/n): ")
				c.inputTrigger <- InputTrigger{Category: PlayAgain}
			case server.MsgTypePlayAgainTimeout:
				fmt.Fprintln(output, "You've been disconnected due to not responding.")
				return nil
			case server.MsgTypeGuessTimeout:
				var guessTimeoutPayload server.GuessTimeoutPayload
				if err := json.Unmarshal(msg.Payload, &guessTimeoutPayload); err != nil {
					log.Println("Error during guess timeout payload unmarshalling:", err)
					return err
				}
				if guessTimeoutPayload.Player.ID == me.ID {
					fmt.Fprintln(output, "Your turn has timed out.")
				} else {
					fmt.Fprintf(output, "Player %s's turn has timed out.\n", guessTimeoutPayload.Player.Nickname)
				}
			}
		case input := <-c.input:
			category := input.Category
			switch category {
			case GuessWord:
				// Check if text is 5
				if len(input.Text) != 5 {
					fmt.Fprintln(output, "Invalid input. Please enter a 5-letter word.")
					fmt.Fprintf(output, "Enter your guess (%d/%d): ", currentRound, maxGuesses)
					c.inputTrigger <- InputTrigger{Category: GuessWord}
					continue
				}
				// Handle guess word input
				guessPayload := server.GuessPayload{
					Word: input.Text,
				}
				payload, err := json.Marshal(guessPayload)
				if err != nil {
					log.Println("Error during guess request marshalling:", err)
					return err
				}
				c.outgoing <- &server.Message{
					Type:    server.MsgTypeGuess,
					Payload: payload,
				}
			case PlayAgain:
				// Handle play again input
				confirmed := input.Text == "y" || input.Text == "Y"
				playAgainPayload, err := json.Marshal(&server.PlayAgainPayload{
					Confirm: confirmed,
				})
				if err != nil {
					log.Println("Error during play again payload marshalling:", err)
					return err
				}
				c.outgoing <- &server.Message{
					Type:    server.MsgTypePlayAgain,
					Payload: playAgainPayload,
				}
				if !confirmed {
					fmt.Fprintln(output, "Thanks for playing!")
					return nil
				}
			}
		case err = <-c.error:
			fmt.Println("Error:", err)
			c.conn.Close()
			return err
		}
	}
}
