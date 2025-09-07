package main

import (
	"fmt"
	"log"

	"github.com/tomlaws/wordle/internal/client"
	"github.com/tomlaws/wordle/internal/controller"
)

func main() {
	// Ask for IP address to connect
	var ipAddress string
	for ipAddress == "" {
		fmt.Print("Enter the server address (localhost:8080): ")
		fmt.Scanln(&ipAddress)
		if ipAddress == "" {
			fmt.Println("No address entered. Defaulting to localhost:8080")
			ipAddress = "localhost:8080"
		}
	}
	var nickname string
	for nickname == "" {
		fmt.Print("Enter your nickname: ")
		fmt.Scanln(&nickname)
		if nickname == "" {
			fmt.Println("Nickname cannot be empty. Please try again.")
		}
	}
	client, err := client.NewClient(ipAddress, nickname)
	if err != nil {
		log.Fatal("Error creating client:", err)
	} else {
		controller.NewController(client)
	}
}
