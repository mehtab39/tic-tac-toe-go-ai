// server/server_test.go

package server_test

import (
	"fmt"
	"go-be-ai/server"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	// Create a request to the /ping endpoint
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call the PingHandler function directly, passing in the ResponseRecorder and Request
	server.PingHandler(rr, req)
	// Check the status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "pong", rr.Body.String())
}

func TestExtractGameID(t *testing.T) {
	// Test valid path
	gameID, ok := server.ExtractGameID("/play/game123")
	assert.True(t, ok, "Expected ok to be true")
	assert.Equal(t, "game123", gameID, "Expected gameID to be game123")
	// Test invalid path
	falsyGameId, ok := server.ExtractGameID("play/")
	assert.Equal(t, "", falsyGameId, "Expected gameID to be ''")
	assert.False(t, ok, "Expected ok to be false")
}

type MockHTTPClient struct {
	Request *http.Request
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Capture the request
	m.Request = req

	// Mock response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{}`)),
	}, nil
}

func TestHandleStartGame(t *testing.T) {
	// Mock environment variables
	os.Setenv("GAME_SERVER", "http://example.com")

	// Create mock HTTP client
	httpClient := &MockHTTPClient{}

	// Call HandleStartGame
	resp, err := server.HandleStartGame("game123", httpClient)

	if err != nil {
		assert.Fail(t, "Got an error")
	}
	defer resp.Body.Close() // Ensure response body is closed
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check request
	assert.Equal(t, "POST", httpClient.Request.Method)                                 // Assertion for request method
	assert.Equal(t, "application/json", httpClient.Request.Header.Get("Content-Type")) // Assertion for Content-Type header

	// Reset environment variables
	os.Unsetenv("GAME_SERVER")
}

func MockWebSocketConnection(gameID string) error {
	// Check if the gameID matches the expected value
	if gameID != "game123" {
		return fmt.Errorf("unexpected gameID: %s", gameID)
	}
	return nil
}

func TestPlayGameHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/play/game123", nil)
	w := httptest.NewRecorder()
	server.PlayGameHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
