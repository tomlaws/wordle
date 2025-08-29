package main

import (
	"log"
	"net/http"

	"github.com/tomlaws/wordle/internal/server"
)

func main() {
	http.HandleFunc("/socket", server.SocketHandler)
	log.Printf("Server starting on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
