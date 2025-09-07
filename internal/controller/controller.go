package controller

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tomlaws/wordle/internal/client"
	"github.com/tomlaws/wordle/internal/game"
	"github.com/tomlaws/wordle/internal/multiplayer"
	"github.com/tomlaws/wordle/internal/protocol"
)

func NewController(client *client.Client) *Controller {
	defer func() {
		client.Stop()
	}()
	protocol := protocol.NewProtocol(multiplayer.PayloadRegistry)
	controller := &Controller{
		input:        make(chan Input),
		inputTrigger: make(chan InputTrigger),
		incoming:     protocol.UnwrapChannel(client.Incoming()),
		outgoing:     protocol.WrapChannel(client.Outgoing()),
		error:        client.Err(),
	}
	go controller.handleInput()
	controller.start()
	return controller
}

func handleInputUnderlying(lines chan interface{}) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines <- s.Text()
	}
	lines <- s.Err()
}

func (c *Controller) handleInput() {
	input := make(chan interface{})
	go handleInputUnderlying(input)
	var tr *InputTrigger = nil
	for {
		select {
		case trigger := <-c.inputTrigger:
			tr = &trigger
		case line := <-input:
			if tr != nil {
				c.input <- Input{Category: tr.Category, Text: line.(string)}
				tr = nil
			}
		}
	}
}

func (c *Controller) start() error {
	defer func() {
		close(c.input)
		close(c.inputTrigger)
	}()
	output := os.Stdout
	var err error
	var me Me
	var maxGuesses int
	var currentRound int
	var isOddPlayer bool
	for {
		select {
		case msg := <-c.incoming:
			switch msg := msg.(type) {
			case *multiplayer.PlayerInfoPayload:
				me = Me{
					ID:       msg.ID,
					Nickname: msg.Nickname,
				}
				fmt.Fprintf(output, "Welcome to Wordle, %s!\n", me.Nickname)
			case *multiplayer.MatchingPayload:
				fmt.Fprintf(output, "Finding opponent...\n")
			case *multiplayer.GameStartPayload:
				// Handle game start
				currentRound = 1
				maxGuesses = msg.MaxGuesses
				isOddPlayer = msg.Player1.ID == me.ID
				var opponent *multiplayer.Player
				if isOddPlayer {
					opponent = msg.Player2
				} else {
					opponent = msg.Player1
				}
				fmt.Fprintf(output, "You are playing against %s\n", opponent.Nickname)
				fmt.Fprintln(output, "Guess the 5-letter word in", maxGuesses, "rounds.")
			case *multiplayer.RoundStartPayload:
				sameRound := msg.Round == currentRound
				currentRound = msg.Round
				timeout := msg.Timeout
				// Handle guess input when it's the player's turn
				if msg.Player.ID == me.ID {
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
			case *multiplayer.InvalidWordPayload:
				if msg.Player.ID == me.ID {
					fmt.Fprintln(output, "Invalid word. Please try again.")
				} else {
					fmt.Fprintf(output, "Opponent guessed an invalid word: %s\n", msg.Word)
				}
			case *multiplayer.FeedbackPayload:
				player := msg.Player
				if player.ID == me.ID {
					fmt.Printf("You guessed: ")
				} else {
					fmt.Printf("Opponent guessed: ")
				}
				currentRound = msg.Round + 1
				// Display feedback to the user
				for _, lr := range msg.Feedback {
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
			case *multiplayer.GameOverPayload:
				if msg.Winner == nil {
					fmt.Fprintln(output, "It's a draw! The correct word was:", msg.Answer)
				} else if msg.Winner.ID == me.ID {
					fmt.Fprintln(output, "Congratulations! You've won!")
				} else {
					fmt.Fprintln(output, "You've lost! The correct word was:", msg.Answer)
				}

				// Ask for a new game
				fmt.Fprint(output, "Do you want to play again? (y/n): ")
				c.inputTrigger <- InputTrigger{Category: PlayAgain}
			case *multiplayer.PlayAgainPayload:
				fmt.Fprintln(output, "You've been disconnected due to not responding.")
				return nil
			case *multiplayer.GuessTimeoutPayload:
				if msg.Player.ID == me.ID {
					fmt.Fprintln(output, "Your turn has timed out.")
				} else {
					fmt.Fprintf(output, "Player %s's turn has timed out.\n", msg.Player.Nickname)
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
				guessPayload := multiplayer.GuessPayload{
					Word: input.Text,
				}
				c.outgoing <- &guessPayload
			case PlayAgain:
				// Handle play again input
				confirmed := input.Text == "y" || input.Text == "Y"
				playAgainPayload := multiplayer.PlayAgainPayload{
					Confirm: confirmed,
				}
				c.outgoing <- &playAgainPayload
				if !confirmed {
					fmt.Fprintln(output, "Thanks for playing!")
					return nil
				}
			}
		case err = <-c.error:
			fmt.Println("Error:", err)
			return err
		}
	}
}
