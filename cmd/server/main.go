package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/tomlaws/wordle/internal/multiplayer"
	"github.com/tomlaws/wordle/internal/server"
)

var Port string = "8080"
var MaxGuesses string = "6"
var ThinkTime string = "60"
var WordListPath string = "assets/words.txt"

func main() {
	var maxGuessesInt int
	if mg, err := strconv.Atoi(MaxGuesses); err == nil {
		maxGuessesInt = mg
	}
	// fatal if max guesses is not greater than or equal to 2 or not divided by 2
	if maxGuessesInt < 2 || maxGuessesInt%2 != 0 {
		log.Fatal("Invalid max guesses. Must be >= 2 and even.")
	}
	var thinkTimeInt int
	if tg, err := strconv.Atoi(ThinkTime); err == nil {
		thinkTimeInt = tg
	}
	if thinkTimeInt < 1 {
		log.Fatal("Invalid think time. Must be >= 1.")
	}
	lobby := multiplayer.NewLobby(WordListPath, maxGuessesInt, thinkTimeInt)
	handler := server.NewServer(
		func(client *server.Client) {
			lobby.NewPlayer(client)
		},
	)
	http.HandleFunc("/socket", handler)
	log.Printf("Server starting on %s", ":"+Port)
	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
