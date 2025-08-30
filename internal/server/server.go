package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tomlaws/wordle/internal/game"
)

var Upgrader = websocket.Upgrader{}
var queue = make(chan *Player, 100) // concurrency-safe queue
var wordList *game.WordList

func startGame(p1, p2 *Player) {
	if p1.connected.Load() && p2.connected.Load() {
		log.Printf("Starting game between %s and %s", p1.Nickname, p2.Nickname)
	} else {
		if p1.connected.Load() {
			log.Printf("Player %s is added back to the queue\n", p1.Nickname)
			queue <- p1
		} else {
			log.Printf("Player %s disconnected before match\n", p1.Nickname)
		}
		if p2.connected.Load() {
			queue <- p2
			log.Printf("Player %s is added back to the queue\n", p2.Nickname)
		} else {
			log.Printf("Player %s disconnected before match\n", p2.Nickname)
		}
		return
	}
	log.Printf("Starting game between %s and %s", p1.Nickname, p2.Nickname)
	// Select random player to start
	gameStartPayload := GameStartPayload{
		MaxAttempts: 12,
	}
	if rand.Intn(2) == 0 {
		gameStartPayload.Player1 = p1
		gameStartPayload.Player2 = p2
	} else {
		gameStartPayload.Player1 = p2
		gameStartPayload.Player2 = p1
	}
	gameStartPayloadBytes, _ := json.Marshal(gameStartPayload)

	// add message to p1's outgoing
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

	for round <= 12 && g.State == game.InProgress {
		msg := <-currentPlayer.incoming
		switch msg.Type {
		case MsgTypeTyping:
			var typingPayload TypingPayload
			if err := json.Unmarshal(msg.Data, &typingPayload); err != nil {
				break
			}
			// Send typing notification to the other player
			if currentPlayer == p1 {
				p2.outgoing <- &Message{
					Type: MsgTypeTyping,
					Data: msg.Data,
				}
			} else {
				p1.outgoing <- &Message{
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
				break
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
}

func matchPlayer() {
	for {
		p1 := <-queue
		p2 := <-queue
		go startGame(p1, p2)
		// Sleep briefly to avoid busy waiting
		time.Sleep(100 * time.Millisecond)
	}
}

func handleRead(player *Player) {
	player.incoming = make(chan *Message)
	defer func() {
		player.connected.Store(false)
		player.conn.Close()
	}()
	for {
		var msg Message
		if err := player.conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message from player %s: %v", player.Nickname, err)
			break
		}
		player.incoming <- &msg
		// Handle incoming messages
		log.Printf("Received message from player %s: %v", player.Nickname, msg.Type)
	}
}

func handleWrite(player *Player) {
	player.outgoing = make(chan *Message)
	defer func() {
		player.connected.Store(false)
		player.conn.Close()
	}()
	for {
		msg := <-player.outgoing
		if err := player.conn.WriteJSON(msg); err != nil {
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
		conn:      conn,
		ID:        uuid.New().String(),
		Nickname:  nickname,
		connected: atomic.Bool{},
	}
	player.connected.Store(true)
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
	queue <- player
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
