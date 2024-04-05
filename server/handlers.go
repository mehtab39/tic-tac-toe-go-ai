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
	gameID, ok := ExtractGameIDFromURL(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Received a game request for game ID:", gameID)

	_, err := handleStartGame(gameID)

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

func handleStartGame(gameID string) (*http.Response, error) {
	// Make a POST request to the other server's API
	fmt.Printf("Start game request from %s \n", GetUserId())
	url := fmt.Sprintf("%s/api/v1/game/start/%s", os.Getenv("GAME_SERVER"), gameID)
	payload := strings.NewReader(fmt.Sprintf(`{"userId": "%s"}`, GetUserId()))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	return resp, nil
}
