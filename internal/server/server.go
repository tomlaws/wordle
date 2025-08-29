package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
	"github.com/tomlaws/wordle/pkg/utils"
)

var Upgrader = websocket.Upgrader{}

func handleConnection(conn *websocket.Conn) {
	var wordlist, err = game.NewWordList(path.Join(utils.Root, "assets", "words.txt"))
	if err != nil {
		log.Println("Error loading word list:", err)
		return
	}

	defer conn.Close()
	for {
		answer := wordlist.RandomWord()
		maxAttempts := 6
		g := game.NewGame(answer, maxAttempts)
		log.Printf("New game started with answer: %s", answer)
		// Send game_start to client
		if err := conn.WriteJSON(Message{
			Type: MsgTypeGameStart,
			Data: json.RawMessage(`{"max_attempts":` + fmt.Sprintf("%d", maxAttempts) + `}`),
		}); err != nil {
			log.Println("Error during game start message sending:", err)
			return
		}

		for g.State == game.InProgress {
			var message Message
			if err := conn.ReadJSON(&message); err != nil {
				log.Println("Error during message reading:", err)
				return
			}

			// Handle message
			switch message.Type {
			case MsgTypeGuess:
				var guess GuessRequest
				if err := json.Unmarshal(message.Data, &guess); err != nil {
					log.Println("Error during guess message unmarshalling:", err)
					return
				}
				// Check if word is valid
				if !wordlist.IsValidWord(guess.Word) {
					// Send invalid_word message
					if err := conn.WriteJSON(Message{
						Type: MsgTypeInvalidWord,
						Data: json.RawMessage(fmt.Sprintf(`{"word":"%s"}`, guess.Word)),
					}); err != nil {
						log.Println("Error during invalid word message sending:", err)
						return
					}
					continue
				}
				// Process guess
				result, err := g.MakeGuess(guess.Word)
				if err != nil {
					log.Println("Error during guess processing:", err)
					return
				}
				// Send result back to client
				var feedbackResponse FeedbackResponse
				feedbackResponse.Feedback = result
				feedbackResponse.Round = len(g.Attempts)
				data, err := json.Marshal(feedbackResponse)
				if err != nil {
					log.Println("Error during feedback marshalling:", err)
					return
				}
				response := Message{
					Type: MsgTypeFeedback,
					Data: data,
				}
				if err := conn.WriteJSON(response); err != nil {
					log.Println("Error during response writing:", err)
					return
				}
				log.Printf("Sent feedback: %+v", feedbackResponse)
			}
		}
		// Check if player would like to play again
		var message Message
		for message.Type != MsgTypeConfirmPlay {
			if err := conn.ReadJSON(&message); err != nil {
				log.Println("Error during confirm play message reading:", err)
				return
			}
		}
		var confirmPlayPayload ConfirmPlayPayload
		if err := json.Unmarshal(message.Data, &confirmPlayPayload); err != nil {
			log.Println("Error during confirm play message unmarshalling:", err)
			return
		}
		if !confirmPlayPayload.Confirm {
			return
		}
	}
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}

	handleConnection(conn)
}
