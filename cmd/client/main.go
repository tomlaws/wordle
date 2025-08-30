package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tomlaws/wordle/internal/client"
)

func main() {
	// Ask for IP address to connect
	var ipAddress string
	fmt.Print("Enter server IP address: ")
	fmt.Scan(&ipAddress)

	var nickname string
	fmt.Print("Enter your nickname: ")
	fmt.Scan(&nickname)

	client, err := client.New(ipAddress, nickname)
	if err != nil {
		log.Fatal("Error creating client:", err)
	}

	if err := client.Start(os.Stdin, os.Stdout); err != nil {
		log.Fatal("Error starting client:", err)
	}
}
