package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
	"github.com/tomlaws/wordle/internal/server"
)

type Client struct {
	conn *websocket.Conn
}

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
	return &Client{conn: conn}, nil
}

func (c *Client) Start(input io.Reader, output io.Writer) error {
	defer c.conn.Close()

	var message server.Message
	if err := c.conn.ReadJSON(&message); err != nil {
		log.Println("Error during message reading:", err)
		return err
	}
	// Read player info
	var playerInfoPayload server.Player
	if err := json.Unmarshal(message.Data, &playerInfoPayload); err != nil {
		log.Println("Error during player info payload unmarshalling:", err)
		return err
	}
	fmt.Fprintf(output, "Welcome to Wordle, %s!\n", playerInfoPayload.Nickname)

	var maxAttempts int
	var currentAttempt int
	var isOddPlayer bool
	for {
		if err := c.conn.ReadJSON(&message); err != nil {
			log.Println("Error during message reading:", err)
			return err
		}
		// log.Printf("Received message: %+v", message)
		switch message.Type {
		case server.MsgTypeGameStart:
			// Handle game start
			var gameStartPayload server.GameStartPayload
			if err := json.Unmarshal(message.Data, &gameStartPayload); err != nil {
				log.Println("Error during game start payload unmarshalling:", err)
				return err
			}
			// Start the game with the received payload
			currentAttempt = 0
			maxAttempts = gameStartPayload.MaxAttempts
			isOddPlayer = gameStartPayload.Player1.ID == playerInfoPayload.ID
			fmt.Fprintln(output, "Guess the 5-letter word in", maxAttempts, "attempts.")
		case server.MsgTypeInvalidWord:
			fmt.Fprintln(output, "Invalid word. Please try again.")
		case server.MsgTypeFeedback:
			var feedbackResponse server.FeedbackResponse
			if err := json.Unmarshal(message.Data, &feedbackResponse); err != nil {
				log.Println("Error during feedback payload unmarshalling:", err)
				return err
			}
			currentAttempt = feedbackResponse.Round
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
			if err := json.Unmarshal(message.Data, &gameOverPayload); err != nil {
				log.Println("Error during game over payload unmarshalling:", err)
				return err
			}
			if gameOverPayload.Winner == nil {
				fmt.Fprintln(output, "It's a draw! The correct word was:", gameOverPayload.Answer)
			} else if gameOverPayload.Winner.ID == playerInfoPayload.ID {
				fmt.Fprintln(output, "Congratulations! You've won!")
			} else {
				fmt.Fprintln(output, "Game over! The correct word was:", gameOverPayload.Answer)
			}

			// Ask for a new game
			fmt.Fprint(output, "Do you want to play again? (y/n): ")
			var playAgain string
			if _, err := fmt.Fscan(input, &playAgain); err != nil {
				log.Println("Error during input reading:", err)
				return err
			}

			confirmPlayPayload := server.ConfirmPlayPayload{
				Confirm: playAgain == "y" || playAgain == "Y",
			}
			var confirmPlayPayloadJson []byte
			confirmPlayPayloadJson, err := json.Marshal(confirmPlayPayload)
			if err != nil {
				log.Println("Error during confirm play payload marshalling:", err)
				return err
			}
			if err := c.conn.WriteJSON(server.Message{
				Type: server.MsgTypeConfirmPlay,
				Data: confirmPlayPayloadJson,
			}); err != nil {
				log.Println("Error during new game request sending:", err)
				return err
			}
			if confirmPlayPayload.Confirm {
				continue
			} else {
				// Disconnect
				fmt.Fprintln(output, "Thanks for playing!")
				return nil
			}
		}
		// Handle guess input when it's the player's turn
		if isOddPlayer && currentAttempt%2 == 0 || !isOddPlayer && currentAttempt%2 == 1 {
			fmt.Fprintf(output, "Enter your guess (%d/%d): ", currentAttempt+1, maxAttempts)
			var guess string
			if _, err := fmt.Fscan(input, &guess); err != nil {
				log.Println("Error during input reading:", err)
				return err
			}
			// Send guess to server
			guessRequest := server.GuessRequest{
				Word: guess,
			}
			data, err := json.Marshal(guessRequest)
			if err != nil {
				log.Println("Error during guess request marshalling:", err)
				return err
			}
			if err := c.conn.WriteJSON(server.Message{
				Type: server.MsgTypeGuess,
				Data: data,
			}); err != nil {
				log.Println("Error during guess message sending:", err)
				return err
			}
		}
	}
}
