package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/tomlaws/wordle/internal/server"
)

var MaxGuesses string = "12"
var WordListPath string = "assets/words.txt"

func main() {
	maxGuessesInt := 6
	if mg, err := strconv.Atoi(MaxGuesses); err == nil {
		maxGuessesInt = mg
	}
	// fatal if max guesses is not greater than or equal to 2 or not divided by 2
	if maxGuessesInt < 2 || maxGuessesInt%2 != 0 {
		log.Fatal("Invalid max guesses. Must be >= 2 and even.")
	}

	http.HandleFunc("/socket", server.Init(WordListPath, maxGuessesInt))
	log.Printf("Server starting on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
