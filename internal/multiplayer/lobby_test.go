package multiplayer

import (
	"encoding/json"
	"path"
	"testing"

	"github.com/tomlaws/wordle/internal/protocol"
	"github.com/tomlaws/wordle/pkg/utils"
)

type MockClient struct {
	id       func() string
	nickname func() string
	incoming func() chan json.RawMessage
	outgoing func() chan json.RawMessage
	err      func() chan error
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

func (m *MockClient) Err() chan error {
	return m.err()
}

func TestLobby_NewPlayer(t *testing.T) {
	lobby := NewLobby(path.Join(utils.Root, "assets", "words.txt"), 6, 30)
	outgoing := make(chan json.RawMessage)
	mockClient := MockClient{
		id:       func() string { return "player1" },
		nickname: func() string { return "Player One" },
		incoming: func() chan json.RawMessage { return make(chan json.RawMessage) },
		outgoing: func() chan json.RawMessage { return outgoing },
		err:      func() chan error { return make(chan error) },
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
