package main

import (
	"fmt"
	"path"

	"github.com/tomlaws/wordle/internal/game"
)

func main() {
	// game loop
	fmt.Println("Welcome to Wordle!")
	for {
		fmt.Println("Guess the 5-letter word in 6 attempts.")
		wordlist, err := game.NewWordList(path.Join("assets", "words.txt"))
		if err != nil {
			fmt.Printf("Error loading word list: %v\n", err)
			return
		}
		answer := wordlist.RandomWord()
		//println("Debug: The answer is", answer) // For testing purposes

		if answer == "" {
			fmt.Println("Word list is empty. Cannot start the game.")
			return
		}
		maxGuesses := 6
		g := game.NewGame(answer, maxGuesses)
		for g.State == game.InProgress {
			var guess string
			fmt.Printf("Enter your guess (%d/%d): ", len(g.Attempts)+1, maxGuesses)
			fmt.Scanln(&guess)
			if len(guess) != 5 {
				fmt.Println("Please enter a 5-letter word.")
				continue
			}
			if wordlist.IsValidWord(guess) == false {
				fmt.Println("Not a valid word. Please try again.")
				continue
			}

			result, err := g.MakeGuess(guess)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			for _, lr := range result {
				switch lr.MatchType {
				case game.Hit:
					fmt.Printf("[%c] ", lr.Letter)
				case game.Present:
					fmt.Printf("(%c) ", lr.Letter)
				case game.Miss:
					fmt.Printf(" %c  ", lr.Letter)
				}
			}
			if g.State == game.InProgress {
				fmt.Println()
			}
		}
		if g.State == game.Won {
			fmt.Println("\nCongratulations! You've guessed the word!")
		} else if g.State == game.Lost {
			fmt.Printf("\nGame over! The correct word was: %s\n", g.Answer)
		}
		// ask to play again
		var playAgain string
		fmt.Print("Play again? (y/n): ")
		fmt.Scanln(&playAgain)
		if playAgain != "y" && playAgain != "Y" {
			fmt.Println("Thanks for playing!")
			break
		}
	}

}
