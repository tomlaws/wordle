package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/tomlaws/wordle/internal/game"
)

func RunGame(input io.Reader, output io.Writer) {
	fmt.Fprintln(output, "Welcome to Wordle!")
	for {
		fmt.Fprintln(output, "Guess the 5-letter word in 6 attempts.")
		wordlist, err := game.NewWordList(path.Join("assets", "words.txt"))
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
		maxGuesses := 6
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
	RunGame(os.Stdin, os.Stdout)
}
