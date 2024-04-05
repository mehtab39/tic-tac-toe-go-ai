package server

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func CreateWebSocketConnection(gameID string) error {
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
	if IsMyTurn(&gameStateUpdate) && gameStateUpdate.GameState.Status == "ONGOING" {
		time.Sleep(time.Millisecond * 500)
		playTurn(conn, &gameStateUpdate)
	}
}

func playTurn(conn *websocket.Conn, gameStateUpdate *GameStateUpdate) {
	row, col := Play(*gameStateUpdate)
	// Publish move
	move := map[string]interface{}{
		"type":   "tileClick",
		"player": GetUserId(),
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
	currentPlayer := "X"
	if gameStateUpdate.GameState.CurrentPlayer == 0 {
		currentPlayer = "O"
	}
	state := GameState{
		Board:         gameStateUpdate.GameState.Board,
		CurrentPlayer: currentPlayer,
		Maximizing:    true,
	}
	depth := 5
	alpha := math.MinInt64
	beta := math.MaxInt64
	_, bestRow, bestCol := Minimax(state, depth, alpha, beta)

	return bestRow, bestCol
}
