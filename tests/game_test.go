package tests

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestGamePlay(t *testing.T) {
	// Run server
	server := exec.Command("go", "run", "./cmd/server")

	go func() {
		if err := server.Run(); err != nil {
			t.Fatalf("Failed to run server: %v", err)
		}
		// output
		output, err := server.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to get server output: %v", err)
		}
		t.Logf("Server output:\n%s", output)
		log.Printf("Server output:\n%s", output)
	}()

	// Run client
	client := exec.Command("go", "run", "./cmd/client")

	go func() {
		if err := client.Run(); err != nil {
			t.Fatalf("Failed to run client: %v", err)
		}
		// output
		output, err := client.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to get client output: %v", err)
		}
		t.Logf("Client output:\n%s", output)
		log.Printf("Client output:\n%s", output)
	}()

	// Input localhost to client
	client.Stdin = strings.NewReader("localhost\n")

	// Logging
}
