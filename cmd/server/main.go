package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/tomlaws/wordle/internal/server"
)

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
	http.HandleFunc("/socket", server.Init(wordListPath, maxGuesses))
	log.Printf("Server starting on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
