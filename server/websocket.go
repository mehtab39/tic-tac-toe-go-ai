package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketConnector is an interface for WebSocket connection.
type WebSocketConnector interface {
	Connect(gameID string) error
}

// DefaultWebSocketConnector is the default implementation of WebSocketConnector.
type DefaultWebSocketConnector struct{}

// Connect creates a WebSocket connection.
func (d *DefaultWebSocketConnector) Connect(gameID string) error {
	fmt.Printf("Creating ws connection at: %s/game/%s\n", os.Getenv("WS_GAME_SERVER"), gameID)
	wsURL := fmt.Sprintf("%s/game/%s", os.Getenv("WS_GAME_SERVER"), gameID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = subscribeToEvents(conn)

	return err
}

// CreateWebSocketConnection creates a WebSocket connection using the provided WebSocketConnector.
func CreateWebSocketConnection(gameID string) error {
	connector := DefaultWebSocketConnector{}
	return connector.Connect(gameID)
}

func subscribeToEvents(conn *websocket.Conn) error {
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
	if isAITurn(&gameStateUpdate) && gameStateUpdate.GameState.Status == "ONGOING" {
		time.Sleep(time.Millisecond * 500)
		playTurn(conn, &gameStateUpdate)
	}
}

func playTurn(conn *websocket.Conn, gameStateUpdate *GameStateUpdate) {
	row, col := Play(*gameStateUpdate)
	// Publish move
	move := map[string]interface{}{
		"type":   "tileClick",
		"player": GetAIUserId(),
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

func Play(gameStateUpdate GameStateUpdate) (int, int) {
	bestRow, bestCol := GetAIMove(gameStateUpdate.GameState.Board)
	return bestRow, bestCol
}
