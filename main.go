package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

// GameStateUpdate represents the structure of a game state update message
type GameStateUpdate struct {
	Type      string `json:"type"`
	GameState struct {
		ID            string     `json:"id"`
		Board         [][]string `json:"board"`
		CreatorID     string     `json:"creatorId"`
		Players       []string   `json:"players"`
		CurrentPlayer int        `json:"currentPlayer"`
		Status        string     `json:"status"`
		Winner        string     `json:"winner"`
	} `json:"gameState"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type StartGameRequest struct {
	UserID string `json:"userId"`
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/play/", playGameHandler)
	fmt.Println("Server is listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a ping request")
	w.Write([]byte("pong"))
}

func handleStartGame(gameID string) (*http.Response, error) {
	// Make a POST request to the other server's API
	url := fmt.Sprintf("http://localhost:5000/api/v1/game/start/%s", gameID)
	payload := strings.NewReader(fmt.Sprintf(`{"userId": "%s"}`, AIGoUserID))

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

func playGameHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := createWebSocketConnection(gameID); err != nil {
		fmt.Println("Error creating WebSocket connection:", err)
		http.Error(w, "Error creating WebSocket connection", http.StatusInternalServerError)
		return
	}

}

func subscribeToEvents(conn *websocket.Conn) error {
	// Continuously read messages from the WebSocket connection
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		// Handle different message types if needed
		switch messageType {
		case websocket.TextMessage:
			fmt.Println("Received test message:", string(message))
			handleTextMessage(message, conn)
		case websocket.BinaryMessage:
			fmt.Println("Received binary message:", message)
		case websocket.CloseMessage:
			fmt.Println("Received close message")
			return nil // Connection closed, return from the function
		}
	}
}

func createWebSocketConnection(gameID string) error {
	// Establish WebSocket connection
	wsURL := fmt.Sprintf("ws://localhost:5000/game/%s", gameID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Subscribe to events
	err = subscribeToEvents(conn)
	if err != nil {
		return err
	}

	return err
}

func handleTextMessage(message []byte, conn *websocket.Conn) {
	var gameStateUpdate GameStateUpdate
	if err := json.Unmarshal(message, &gameStateUpdate); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if gameStateUpdate.GameState.Status == "FINISHED" {
		// Close the WebSocket connection if the game is finished
		fmt.Println("Game finished. Closing WebSocket connection.")
		conn.Close()
		return
	}
	// Check if it's AI's turn
	if IsMyTurn(&gameStateUpdate) && gameStateUpdate.GameState.Status == "ONGOING" {
		playTurn(conn, &gameStateUpdate)
	}
}

func playTurn(conn *websocket.Conn, gameStateUpdate *GameStateUpdate) {
	row, col := Play(*gameStateUpdate)
	// Publish move
	move := map[string]interface{}{
		"type":   "tileClick",
		"player": AIGoUserID,
		"row":    row,
		"col":    col,
	}
	moveJSON, err := json.Marshal(move)
	if err != nil {
		fmt.Println("Error marshalling move JSON:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, moveJSON); err != nil {
		fmt.Println("Error publishing move:", err)
		return
	}
}
