package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSocketHandler_UpgradeFailure(t *testing.T) {
	// Create a request that will fail to upgrade (not a websocket request)
	req := httptest.NewRequest("GET", "/socket", nil)
	w := httptest.NewRecorder()
	SocketHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest && res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected 400 or 500 status code for non-websocket request, got %d", res.StatusCode)
	}
}

func TestSocketHandler_UpgradeSuccess_GameStart(t *testing.T) {
	// Start a test server with the SocketHandler
	ts := httptest.NewServer(http.HandlerFunc(SocketHandler))
	defer ts.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u, _ := url.Parse(ts.URL)
	u.Scheme = "ws"
	u.Path = "/socket"

	// Connect to the server using gorilla/websocket
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("WebSocket dial failed: %v", err)
	}
	defer ws.Close()
}
