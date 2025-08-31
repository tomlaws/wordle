package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
)

var Upgrader = websocket.Upgrader{}
var queue = make(chan *Player, 100) // concurrency-safe queue
var wordList *game.WordList

func checkPlayAgain(player *Player) bool {
	timeout := time.After(10 * time.Second)
	select {
	case msg := <-player.incoming:
		if msg.Type == MsgTypePlayAgain {
			var confirmPlayPayload ConfirmPlayPayload
			if err := json.Unmarshal(msg.Data, &confirmPlayPayload); err != nil {
				log.Println("Error unmarshalling play again response:", err)
				return false
			}
			if !confirmPlayPayload.Confirm {
				log.Printf("Player %s declined to play again", player.Nickname)
			} else {
				log.Printf("Player %s wants to play again", player.Nickname)
				enqueuePlayer(player)
			}
			return confirmPlayPayload.Confirm
		}
	case <-timeout:
		log.Printf("Disconnected: Player %s did not respond to play again prompt", player.Nickname)
		player.error <- nil
	}
	return false
}

func startGame(p1, p2 *Player) {
	log.Printf("Starting game between %s and %s", p1.Nickname, p2.Nickname)
	// Select random player to start
	gameStartPayload := GameStartPayload{
		MaxAttempts: 12,
	}
	// Player 1 goes first
	if rand.Intn(2) == 0 {
		gameStartPayload.Player1 = p1
		gameStartPayload.Player2 = p2
	} else {
		gameStartPayload.Player1 = p2
		gameStartPayload.Player2 = p1
	}
	gameStartPayloadBytes, _ := json.Marshal(gameStartPayload)

	p1.outgoing <- &Message{
		Type: MsgTypeGameStart,
		Data: gameStartPayloadBytes,
	}
	p2.outgoing <- &Message{
		Type: MsgTypeGameStart,
		Data: gameStartPayloadBytes,
	}

	currentPlayer := gameStartPayload.Player1
	g := game.NewGame(wordList.RandomWord(), 12)
	log.Printf("Game started with answer: %s", g.Answer)
	round := 1
	var winner *Player

	for round <= 12 && g.State == game.InProgress && winner == nil {
		var turnStartPayload TurnStartPayload
		turnStartPayload.Player = currentPlayer
		data, err := json.Marshal(turnStartPayload)
		if err != nil {
			log.Println("Error during turn start payload marshalling:", err)
			return
		}
		p1.outgoing <- &Message{
			Type: MsgTypeTurnStart,
			Data: data,
		}
		p2.outgoing <- &Message{
			Type: MsgTypeTurnStart,
			Data: data,
		}
		select {
		case p1Err := <-p1.error:
			log.Println("Error from player 1:", p1Err)
			winner = p2
		case p2Err := <-p2.error:
			log.Println("Error from player 2:", p2Err)
			winner = p1
		case msg := <-currentPlayer.incoming:
			switch msg.Type {
			case MsgTypeTyping:
				var typingPayload TypingPayload
				if err := json.Unmarshal(msg.Data, &typingPayload); err != nil {
					break
				}
				// Send typing notification to the other player
				if currentPlayer == p1 {
					p1.outgoing <- &Message{
						Type: MsgTypeTyping,
						Data: data,
					}
				} else {
					p2.outgoing <- &Message{
						Type: MsgTypeTyping,
						Data: msg.Data,
					}
				}
			case MsgTypeGuess:
				// Handle guess
				log.Printf("Player %s guessed", currentPlayer.Nickname)
				var guessRequest GuessRequest
				if err := json.Unmarshal(msg.Data, &guessRequest); err != nil {
					break
				}
				// Process the guess
				result, _ := g.MakeGuess(guessRequest.Word)
				if g.State == game.Won {
					winner = currentPlayer
				}
				// Send the feedback to both players
				var feedbackResponse FeedbackResponse
				feedbackResponse.Round = round
				feedbackResponse.Feedback = result
				data, err := json.Marshal(feedbackResponse)
				if err != nil {
					log.Println("Error during feedback marshalling:", err)
					return
				}
				p1.outgoing <- &Message{
					Type: MsgTypeFeedback,
					Data: data,
				}
				p2.outgoing <- &Message{
					Type: MsgTypeFeedback,
					Data: data,
				}
				// Swap players
				round++
				if currentPlayer == p1 {
					currentPlayer = p2
				} else {
					currentPlayer = p1
				}
			}
		}
	}
	// Game over
	var gameOverPayload GameOverPayload
	if winner != nil {
		log.Printf("Player %s wins!", winner.Nickname)
		gameOverPayload.Winner = winner
		gameOverPayload.Answer = g.Answer
	} else {
		log.Printf("Game ended in a draw")
		gameOverPayload.Winner = nil
		gameOverPayload.Answer = g.Answer
	}
	log.Printf("Game over: %+v", gameOverPayload)
	data, err := json.Marshal(gameOverPayload)
	if err != nil {
		log.Println("Error during game over marshalling:", err)
		return
	}
	p1.outgoing <- &Message{
		Type: MsgTypeGameOver,
		Data: data,
	}
	p2.outgoing <- &Message{
		Type: MsgTypeGameOver,
		Data: data,
	}
	go checkPlayAgain(p1)
	go checkPlayAgain(p2)
}

func matchPlayer() {
	for {
		p1 := <-queue
		p2 := <-queue
		go func() {
			timeout := time.After(2 * time.Second)
			select {
			case <-p1.error:
				log.Printf("Player %s has disconnected", p1.Nickname)
				queue <- p2
			case <-p2.error:
				log.Printf("Player %s has disconnected", p2.Nickname)
				queue <- p1
			case <-timeout:
				log.Printf("Start game between %s and %s", p1.Nickname, p2.Nickname)
				startGame(p1, p2)
			}
		}()
		// Sleep briefly to avoid busy waiting
		time.Sleep(100 * time.Millisecond)
	}
}

func enqueuePlayer(player *Player) {
	if data, err := json.Marshal(player); err == nil {
		player.outgoing <- &Message{
			Type: MsgTypeMatching,
			Data: data,
		}
	}

	select {
	case queue <- player:
		log.Printf("Player %s added to queue", player.Nickname)
	default:
		log.Printf("Player %s could not be added to queue", player.Nickname)
	}
}

func handleRead(player *Player) {
	defer func() {
		player.conn.Close()
	}()
	for {
		var msg Message
		if err := player.conn.ReadJSON(&msg); err != nil {
			player.error <- err
			log.Printf("Error reading message from player %s: %v", player.Nickname, err)
			break
		}
		player.incoming <- &msg
		// Handle incoming messages
		log.Printf("Received message from player %s: %v", player.Nickname, msg.Type)
	}
}

func handleWrite(player *Player) {
	defer func() {
		player.conn.Close()
	}()
	for {
		msg := <-player.outgoing
		if err := player.conn.WriteJSON(msg); err != nil {
			player.error <- err
			log.Printf("Error sending message to player %s: %v", player.Nickname, err)
			break
		}
		// Handle outgoing messages
		log.Printf("Sending message to player %s: %s", player.Nickname, msg.Type)
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	nickname := r.URL.Query().Get("nickname")
	player := &Player{
		conn:     conn,
		ID:       uuid.New().String(),
		Nickname: nickname,
		error:    make(chan error),
		incoming: make(chan *Message),
		outgoing: make(chan *Message),
	}
	go handleRead(player)
	go handleWrite(player)

	log.Printf("New player connected: %s (%s)", player.Nickname, player.ID)
	// Send player info to client
	playerInfoBytes, err := json.Marshal(player)
	if err != nil {
		log.Printf("Error marshalling player info: %v", err)
		return
	}
	player.outgoing <- &Message{
		Type: MsgTypePlayerInfo,
		Data: playerInfoBytes,
	}

	// Add the new player to the queue
	enqueuePlayer(player)
}

func Init(wordListPath string, maxAttempts int) func(w http.ResponseWriter, r *http.Request) {
	var err error
	wordList, err = game.NewWordList(wordListPath)
	if err != nil {
		log.Fatalf("Error loading word list: %v", err)
	}
	go matchPlayer()
	return socketHandler
}
