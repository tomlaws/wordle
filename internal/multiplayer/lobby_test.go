package multiplayer

import (
	"encoding/json"
	"path"
	"testing"
	"time"

	"github.com/tomlaws/wordle/internal/protocol"
	"github.com/tomlaws/wordle/pkg/utils"
)

type MockClient struct {
	id       func() string
	nickname func() string
	incoming func() chan json.RawMessage
	outgoing func() chan json.RawMessage
	error    func() chan error
}

func (m *MockClient) ID() string {
	return m.id()
}

func (m *MockClient) Nickname() string {
	return m.nickname()
}

func (m *MockClient) Incoming() chan json.RawMessage {
	return m.incoming()
}

func (m *MockClient) Outgoing() chan json.RawMessage {
	return m.outgoing()
}

func (m *MockClient) Error() chan error {
	return m.error()
}

func TestLobby_NewPlayer(t *testing.T) {
	lobby := NewLobby(path.Join(utils.Root, "assets", "words.txt"), 6, 30)
	outgoing := make(chan json.RawMessage)
	mockClient := MockClient{
		id:       func() string { return "player1" },
		nickname: func() string { return "Player One" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing },
		error:    func() chan error { return make(chan error) },
	}
	go func() {
		lobby.NewPlayer(&mockClient)
	}()
	msg := <-mockClient.outgoing()
	var message protocol.Message
	if err := json.Unmarshal(msg, &message); err != nil {
		t.Fatalf("Failed to unmarshal Message: %v", err)
	}
	if message.Type != "player_info" {
		t.Errorf("Expected message type 'player_info', got '%s'", message.Type)
	}
	var payload PlayerInfoPayload
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		t.Fatalf("Failed to unmarshal PlayerInfoPayload: %v", err)
	}
	if payload.ID != "player1" {
		t.Errorf("Expected player ID 'player1', got '%s'", payload.ID)
	}
}

func TestLobby_RemovePlayer(t *testing.T) {
	lobby := NewLobby(path.Join(utils.Root, "assets", "words.txt"), 6, 30)
	outgoing := make(chan json.RawMessage)
	mockClient := MockClient{
		id:       func() string { return "player1" },
		nickname: func() string { return "Player One" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing },
		error:    func() chan error { return make(chan error) },
	}
	go func() {
		<-mockClient.outgoing()
	}()
	player := lobby.NewPlayer(&mockClient)
	lobby.RemovePlayer(player)
	_, ok := <-player.outgoing
	if ok {
		t.Errorf("Expected player outgoing channel to be closed")
	}
	_, ok = <-player.incoming
	if ok {
		t.Errorf("Expected player incoming channel to be closed")
	}
}

func TestLobby_AddPlayerToQueue(t *testing.T) {
	lobby := NewLobby(path.Join(utils.Root, "assets", "words.txt"), 6, 30)
	outgoing1 := make(chan json.RawMessage)
	mockClient1 := MockClient{
		id:       func() string { return "player1" },
		nickname: func() string { return "Player One" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing1 },
		error:    func() chan error { return make(chan error) },
	}
	go func() {
		<-mockClient1.outgoing()
	}()
	lobby.NewPlayer(&mockClient1)
	if len(lobby.queue) != 1 {
		t.Errorf("Expected queue length 1 after adding first player, got %d", len(lobby.queue))
	}
	outgoing2 := make(chan json.RawMessage)
	mockClient2 := MockClient{
		id:       func() string { return "player2" },
		nickname: func() string { return "Player Two" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing2 },
		error:    func() chan error { return make(chan error) },
	}
	go func() {
		<-mockClient2.outgoing()
	}()
	lobby.NewPlayer(&mockClient2)
	if len(lobby.queue) != 2 {
		t.Errorf("Expected queue length 2 after adding second player, got %d", len(lobby.queue))
	}
}

func TestLobby_SkipDisconnectedPlayer(t *testing.T) {
	lobby := NewLobby(path.Join(utils.Root, "assets", "words.txt"), 6, 5)
	outgoing1 := make(chan json.RawMessage)
	error1 := make(chan error, 1)
	mockClient1 := MockClient{
		id:       func() string { return "player1" },
		nickname: func() string { return "Player One" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing1 },
		error:    func() chan error { return error1 },
	}
	go func() {
		<-mockClient1.outgoing()
	}()
	lobby.NewPlayer(&mockClient1)
	if len(lobby.queue) != 1 {
		t.Errorf("Expected queue length 1 after adding first player, got %d", len(lobby.queue))
	}
	outgoing2 := make(chan json.RawMessage)
	mockClient2 := MockClient{
		id:       func() string { return "player2" },
		nickname: func() string { return "Player Two" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing2 },
		error:    func() chan error { return make(chan error) },
	}
	go func() {
		<-mockClient2.outgoing()
	}()
	lobby.NewPlayer(&mockClient2)
	if len(lobby.queue) != 2 {
		t.Errorf("Expected queue length 2 after adding second player, got %d", len(lobby.queue))
	}
	// Simulate player 1 disconnecting
	error1 <- nil
	// Wait to allow matching to process
	time.Sleep(100 * time.Millisecond)
	if len(lobby.queue) != 1 {
		t.Errorf("Expected queue length 1 after player 1 disconnected, got %d", len(lobby.queue))
	}
}
