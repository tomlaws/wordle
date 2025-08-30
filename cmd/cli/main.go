package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/tomlaws/wordle/internal/game"
)

func RunGame(input io.Reader, output io.Writer, wordListPath string, maxGuesses int) {
	fmt.Fprintln(output, "Welcome to Wordle!")
	for {
		fmt.Fprintf(output, "Guess the 5-letter word in %d attempts.\n", maxGuesses)
		wordlist, err := game.NewWordList(wordListPath)
		if err != nil {
			fmt.Fprintf(output, "Error loading word list: %v\n", err)
			return
		}
		answer := wordlist.RandomWord()
		//println("Debug: The answer is", answer) // For testing purposes

		if answer == "" {
			fmt.Fprintln(output, "Word list is empty. Cannot start the game.")
			return
		}
		g := game.NewGame(answer, maxGuesses)
		for g.State == game.InProgress {
			var guess string
			fmt.Fprintf(output, "Enter your guess (%d/%d): ", len(g.Attempts)+1, maxGuesses)
			fmt.Fscanln(input, &guess)
			if len(guess) != 5 {
				fmt.Fprintln(output, "Please enter a 5-letter word.")
				continue
			}
			if wordlist.IsValidWord(guess) == false {
				fmt.Fprintln(output, "Not a valid word. Please try again.")
				continue
			}

			result, err := g.MakeGuess(guess)
			if err != nil {
				fmt.Fprintf(output, "Error: %v\n", err)
				continue
			}
			for _, lr := range result {
				switch lr.MatchType {
				case game.Hit:
					fmt.Fprintf(output, "[%c] ", lr.Letter)
				case game.Present:
					fmt.Fprintf(output, "(%c) ", lr.Letter)
				case game.Miss:
					fmt.Fprintf(output, " %c  ", lr.Letter)
				}
			}
			if g.State == game.InProgress {
				fmt.Fprintln(output)
			}
		}
		if g.State == game.Won {
			fmt.Fprintln(output, "\nCongratulations! You've guessed the word!")
		} else if g.State == game.Lost {
			fmt.Fprintf(output, "\nGame over! The correct word was: %s\n", g.Answer)
		}
		// ask to play again
		var playAgain string
		fmt.Fprint(output, "Play again? (y/n): ")
		fmt.Fscanln(input, &playAgain)
		if playAgain != "y" && playAgain != "Y" {
			fmt.Fprintln(output, "Thanks for playing!")
			break
		}
	}

}

func main() {
	// Default settings
	wordListPath := "assets/words.txt"
	maxGuesses := 6

	// Override default values with .env settings
	err := godotenv.Load()
	if err != nil {
		log.Print("Cannot load .env file. Using default settings.")
	} else {
		if mg, err := strconv.Atoi(os.Getenv("MAX_GUESSES")); mg != 0 && err == nil {
			maxGuesses = mg
		} else {
			log.Print("Invalid MAX_GUESSES value. Using default.")
		}
		if wlPath := os.Getenv("WORDLIST_PATH"); wlPath != "" {
			wordListPath = wlPath
		} else {
			log.Print("WORDLIST_PATH not set. Using default.")
		}
	}

	RunGame(os.Stdin, os.Stdout, wordListPath, maxGuesses)
}
