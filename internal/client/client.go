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

func handleInput(client *Client) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		select {
		case trigger := <-client.inputTrigger:
			client.input <- Input{Category: trigger.Category, Text: text}
		default:
			// Ignore input if not triggered
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
	log.Println("Disconnected from server.")
	os.Exit(0)
}

func (c *Client) Start(input io.Reader, output io.Writer) error {
	defer c.Stop()
	go handleInput(c)
	go handleRead(c)
	go handleWrite(c)
	var err error
	var player server.Player
	var maxAttempts int
	var currentAttempt int
	var isOddPlayer bool
	for {
		select {
		case msg := <-c.incoming:
			switch msg.Type {
			case server.MsgTypePlayerInfo:
				if err := json.Unmarshal(msg.Data, &player); err != nil {
					log.Println("Error during player info payload unmarshalling:", err)
					return err
				}
				fmt.Fprintf(output, "Welcome to Wordle, %s!\n", player.Nickname)
			case server.MsgTypeMatching:
				fmt.Fprintf(output, "Finding opponent...\n")
			case server.MsgTypeGameStart:
				// Handle game start
				var gameStartPayload server.GameStartPayload
				if err := json.Unmarshal(msg.Data, &gameStartPayload); err != nil {
					log.Println("Error during game start payload unmarshalling:", err)
					return err
				}
				// Start the game with the received payload
				currentAttempt = 1
				maxAttempts = gameStartPayload.MaxAttempts
				isOddPlayer = gameStartPayload.Player1.ID == player.ID
				fmt.Fprintln(output, "Guess the 5-letter word in", maxAttempts, "attempts.")
			case server.MsgTypeTurnStart:
				var turnStartPayload server.TurnStartPayload
				if err := json.Unmarshal(msg.Data, &turnStartPayload); err != nil {
					log.Println("Error during turn start payload unmarshalling:", err)
					return err
				}
				// Handle guess input when it's the player's turn
				if turnStartPayload.Player.ID == player.ID {
					fmt.Fprintf(output, "Enter your guess (%d/%d): ", currentAttempt, maxAttempts)
					c.inputTrigger <- InputTrigger{Category: GuessWord}
				} else {
					// Wait for opponent's guess
					log.Println("Waiting for opponent's guess...")
				}
			case server.MsgTypeInvalidWord:
				fmt.Fprintln(output, "Invalid word. Please try again.")
			case server.MsgTypeFeedback:
				var feedbackResponse server.FeedbackResponse
				if err := json.Unmarshal(msg.Data, &feedbackResponse); err != nil {
					log.Println("Error during feedback payload unmarshalling:", err)
					return err
				}
				if isOddPlayer && feedbackResponse.Round%2 == 0 || !isOddPlayer && feedbackResponse.Round%2 == 1 {
					fmt.Printf("Opponent guessed: ")
				} else {
					fmt.Printf("You guessed: ")
				}
				currentAttempt = feedbackResponse.Round + 1
				// Display feedback to the user
				for _, lr := range feedbackResponse.Feedback {
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
				if err := json.Unmarshal(msg.Data, &gameOverPayload); err != nil {
					log.Println("Error during game over payload unmarshalling:", err)
					return err
				}
				if gameOverPayload.Winner == nil {
					fmt.Fprintln(output, "It's a draw! The correct word was:", gameOverPayload.Answer)
				} else if gameOverPayload.Winner.ID == player.ID {
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
			}
		case input := <-c.input:
			category := input.Category
			switch category {
			case GuessWord:
				// Check if text is 5
				if len(input.Text) != 5 {
					fmt.Fprintln(output, "Invalid input. Please enter a 5-letter word.")
					fmt.Fprintf(output, "Enter your guess (%d/%d): ", currentAttempt, maxAttempts)
					c.inputTrigger <- InputTrigger{Category: GuessWord}
					continue
				}
				// Handle guess word input
				guessRequest := server.GuessRequest{
					Word: input.Text,
				}
				data, err := json.Marshal(guessRequest)
				if err != nil {
					log.Println("Error during guess request marshalling:", err)
					return err
				}
				c.outgoing <- &server.Message{
					Type: server.MsgTypeGuess,
					Data: data,
				}
			case PlayAgain:
				// Handle play again input
				confirmPlayPayload := server.ConfirmPlayPayload{
					Confirm: input.Text == "y" || input.Text == "Y",
				}
				confirmPlayPayloadJson, err := json.Marshal(confirmPlayPayload)
				if err != nil {
					log.Println("Error during confirm play payload marshalling:", err)
					return err
				}
				c.outgoing <- &server.Message{
					Type: server.MsgTypePlayAgain,
					Data: confirmPlayPayloadJson,
				}
				if confirmPlayPayload.Confirm {
					continue
				} else {
					// Disconnect
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
