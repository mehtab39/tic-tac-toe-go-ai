package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a ping request")
	w.Write([]byte("pong"))
}

func PlayGameHandler(w http.ResponseWriter, r *http.Request) {
	gameID, ok := ExtractGameID(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Received a game request for game ID:", gameID)

	_, err := HandleStartGame(gameID, &http.Client{})

	if err != nil {
		return
	}

	// Create WebSocket connection
	if err := CreateWebSocketConnection(gameID); err != nil {
		fmt.Println("Error creating WebSocket connection:", err)
		http.Error(w, "Error creating WebSocket connection", http.StatusInternalServerError)
		return
	}

}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RealHTTPClient implements the HTTPClient interface
type RealHTTPClient struct{}

// Do implements the Do method of the HTTPClient interface
func (c *RealHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

// HandleStartGame makes a POST request to start a game
func HandleStartGame(gameID string, client HTTPClient) (*http.Response, error) {
	// Make a POST request to the other server's API
	fmt.Printf("Start game request from %s \n", GetAIUserId())
	url := fmt.Sprintf("%s/api/v1/game/start/%s", os.Getenv("GAME_SERVER"), gameID)
	payload := strings.NewReader(fmt.Sprintf(`{"userId": "%s"}`, GetAIUserId()))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func ExtractGameID(path string) (string, bool) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return "", false
	}
	return parts[2], true
}
