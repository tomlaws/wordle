package multiplayer

import (
	"log"
	"math/rand"
	"time"

	"github.com/tomlaws/wordle/internal/game"
	"github.com/tomlaws/wordle/internal/protocol"
)

func NewLobby(wordListPath string, maxGuesses int, thinkTime int) *Lobby {
	wordList, err := game.NewWordList(wordListPath)
	if err != nil {
		log.Fatal("Error loading word list:", err)
	}
	lobby := &Lobby{
		wordList:   wordList,
		maxGuesses: maxGuesses,
		thinkTime:  thinkTime,
		queue:      make(chan *Player, 100),
	}
	go lobby.startMatchingPlayer()
	return lobby
}

func (l *Lobby) NewPlayer(client Client) *Player {
	log.Printf("New player connected: %s", client.Nickname())
	protocol := protocol.NewProtocol(PayloadRegistry)
	player := &Player{
		ID:       client.ID(),
		Nickname: client.Nickname(),
		incoming: protocol.UnwrapChannel(client.Incoming()),
		outgoing: protocol.WrapChannel(client.Outgoing()),
		err:      client.Err(),
	}
	// Welcome
	player.outgoing <- &PlayerInfoPayload{
		ID:       player.ID,
		Nickname: player.Nickname,
	}
	l.addPlayer(player)
	return player
}

func (l *Lobby) RemovePlayer(player *Player) {
	log.Printf("Removing player: %s", player.Nickname)
	close(player.incoming)
	close(player.outgoing)
}

func (l *Lobby) startGame(p1, p2 *Player) {
	// Select random player to start
	gameStartPayload := GameStartPayload{
		MaxGuesses: l.maxGuesses,
	}
	// Player 1 goes first
	if rand.Intn(2) == 0 {
		gameStartPayload.Player1 = p1
		gameStartPayload.Player2 = p2
	} else {
		gameStartPayload.Player1 = p2
		gameStartPayload.Player2 = p1
	}
	p1.outgoing <- &gameStartPayload
	p2.outgoing <- &gameStartPayload

	currentPlayer := gameStartPayload.Player1
	g := game.NewGame(l.wordList.RandomWord(), gameStartPayload.MaxGuesses)
	log.Printf("Game started with answer: %s", g.Answer)
	round := 1
	timeout := time.Duration(l.thinkTime) * time.Second
	var winner *Player

	sendRoundStart := func(player *Player, round int) {
		var roundStartPayload RoundStartPayload
		roundStartPayload.Player = player
		roundStartPayload.Round = round
		roundStartPayload.Deadline = time.Now().Add(timeout)
		p1.outgoing <- &roundStartPayload
		p2.outgoing <- &roundStartPayload
	}

	sendRoundStart(currentPlayer, round)

	for round <= 12 && g.State == game.InProgress && winner == nil {
		roundTimeout := time.After(timeout)

		select {
		case p1Err := <-p1.err:
			log.Println("Error from player 1:", p1Err)
			winner = p2
		case p2Err := <-p2.err:
			log.Println("Error from player 2:", p2Err)
			winner = p1
		case <-roundTimeout:
			log.Println("Guess timeout for player:", currentPlayer.Nickname)
			// Send timeout message
			var guessTimeoutPayload GuessTimeoutPayload
			guessTimeoutPayload.Player = currentPlayer
			guessTimeoutPayload.Round = round
			p1.outgoing <- &guessTimeoutPayload
			p2.outgoing <- &guessTimeoutPayload
			// Swap players and increment round
			round++
			if round <= 12 {
				if currentPlayer == p1 {
					currentPlayer = p2
				} else {
					currentPlayer = p1
				}
				sendRoundStart(currentPlayer, round)
			}
		case rawMsg := <-currentPlayer.incoming:
			switch msg := rawMsg.(type) {
			case *TypingPayload:
				log.Printf("Player %s is typing: %s", currentPlayer.Nickname, msg.Word)
				// Send to the other player
				if currentPlayer == p1 {
					msg.Player = p1
					p2.outgoing <- msg
				} else {
					msg.Player = p2
					p1.outgoing <- msg
				}
			case *GuessPayload:
				// Handle guess
				log.Printf("Player %s guessed: %s", currentPlayer.Nickname, msg.Word)
				// Validate the word
				if !l.wordList.IsValidWord(msg.Word) {
					log.Printf("Invalid word guessed")
					var invalidWordPayload InvalidWordPayload
					invalidWordPayload.Player = currentPlayer
					invalidWordPayload.Round = round
					invalidWordPayload.Word = msg.Word
					p1.outgoing <- &invalidWordPayload
					p2.outgoing <- &invalidWordPayload
					continue
				}
				// Process the guess
				result, _ := g.MakeGuess(msg.Word)
				if g.State == game.Won {
					winner = currentPlayer
				}
				// Send the feedback to both players
				var feedbackPayload FeedbackPayload
				feedbackPayload.Player = currentPlayer
				feedbackPayload.Round = round
				feedbackPayload.Feedback = result
				p1.outgoing <- &feedbackPayload
				p2.outgoing <- &feedbackPayload
				// Swap players and increment round
				round++
				if round <= 12 && winner == nil && g.State == game.InProgress {
					if currentPlayer == p1 {
						currentPlayer = p2
					} else {
						currentPlayer = p1
					}
					sendRoundStart(currentPlayer, round)
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
	p1.outgoing <- &gameOverPayload
	p2.outgoing <- &gameOverPayload
	go l.checkPlayAgain(p1)
	go l.checkPlayAgain(p2)
}

func (l *Lobby) checkPlayAgain(player *Player) bool {
	rawMsg := <-player.incoming
	switch msg := rawMsg.(type) {
	case *PlayAgainPayload:
		if !msg.Confirm {
			log.Printf("Player %s declined to play again", player.Nickname)
		} else {
			log.Printf("Player %s wants to play again", player.Nickname)
			l.addPlayer(player)
		}
		return msg.Confirm
	}
	return false
}

func (l *Lobby) startMatchingPlayer() {
	for {
		p1 := <-l.queue
		p2 := <-l.queue
		go func() {
			timeout := time.After(2 * time.Second)
			select {
			case <-p1.err:
				log.Printf("Player %s has disconnected", p1.Nickname)
				l.queue <- p2
			case <-p2.err:
				log.Printf("Player %s has disconnected", p2.Nickname)
				l.queue <- p1
			case <-timeout:
				log.Printf("Starting game between %s and %s", p1.Nickname, p2.Nickname)
				l.startGame(p1, p2)
			}
		}()
		// Sleep briefly to avoid busy waiting
		time.Sleep(100 * time.Millisecond)
	}
}

func (l *Lobby) addPlayer(player *Player) {
	player.outgoing <- &MatchingPayload{}
	select {
	case l.queue <- player:
		log.Printf("Player %s added to queue", player.Nickname)
	default:
		log.Printf("Player %s could not be added to queue", player.Nickname)
	}
}
