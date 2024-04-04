package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math"
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
	payload := strings.NewReader(`{"userId": "ai-go"}`)

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
	gameID, ok := extractGameIDFromURL(r.URL.Path)
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

func extractGameIDFromURL(path string) (string, bool) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return "", false
	}
	return parts[2], true
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

	// Check if it's AI's turn
	if gameStateUpdate.GameState.Players[gameStateUpdate.GameState.CurrentPlayer] == "ai-go" {
		// Check if the game is ongoing
		if gameStateUpdate.GameState.Status == "ONGOING" {
			currentPlayer := "X"
			if gameStateUpdate.GameState.CurrentPlayer == 0 {
				currentPlayer = "O"
			}
			state := GameState{
				Board:         gameStateUpdate.GameState.Board,
				CurrentPlayer: currentPlayer,
				Maximizing:    true, // Maximizing player is X
			}
			depth := 5 // or 6, 7 depending on your preference
			alpha := math.MinInt64
			beta := math.MaxInt64
			_, bestRow, bestCol := Minimax(state, depth, alpha, beta)

			// Publish move
			move := map[string]interface{}{
				"type":   "tileClick",
				"player": "ai-go",
				"row":    bestRow,
				"col":    bestCol,
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
		} else if gameStateUpdate.GameState.Status == "FINISHED" {
			// Close the WebSocket connection if the game is finished
			fmt.Println("Game finished. Closing WebSocket connection.")
			conn.Close()
		}
	}
}

// Represents the state of the Tic-Tac-Toe game
type GameState struct {
	Board         [][]string // 3x3 board
	CurrentPlayer string     // Current player ('X' or 'O')
	Maximizing    bool       // Indicates if it's the maximizing player's turn
}

// Evaluate function to determine the score of the current state
func (state GameState) Evaluate() int {
	// Winning conditions
	winningLines := [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}
	for _, line := range winningLines {
		if state.Board[line[0]/3][line[0]%3] == state.Board[line[1]/3][line[1]%3] &&
			state.Board[line[1]/3][line[1]%3] == state.Board[line[2]/3][line[2]%3] {
			if state.Board[line[0]/3][line[0]%3] == "O" {
				return -1
			} else if state.Board[line[0]/3][line[0]%3] == "X" {
				return 1
			}
		}
	}

	// Draw condition
	for _, row := range state.Board {
		for _, cell := range row {
			if cell == "" {
				// Game not finished yet
				return 0
			}
		}
	}

	// Game is a draw
	return 0
}

// Minimax function to determine the best move with alpha-beta pruning
func Minimax(state GameState, depth int, alpha, beta int) (int, int, int) {
	// Base case: check if the game is over or max depth reached
	score := state.Evaluate()
	if score != 0 || depth == 0 {
		// Return the score along with invalid row and column
		return score, -1, -1
	}

	// Initialize bestScore to the worst possible score for the maximizing player
	bestScore := math.MinInt64
	bestRow := -1
	bestCol := -1

	// Loop through all possible moves
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if state.Board[i][j] == "" {
				// Make the move
				state.Board[i][j] = state.CurrentPlayer

				// Switch player
				if state.CurrentPlayer == "O" {
					state.CurrentPlayer = "X"
				} else {
					state.CurrentPlayer = "O"
				}

				// Recursively call Minimax for the next state with reduced depth and updated alpha, beta
				score, _, _ := Minimax(state, depth-1, alpha, beta)

				// Undo the move
				state.Board[i][j] = ""
				if state.CurrentPlayer == "O" {
					state.CurrentPlayer = "X"
				} else {
					state.CurrentPlayer = "O"
				}

				// Update bestScore and bestMove based on the current player
				if state.Maximizing {
					if score > bestScore {
						bestScore = score
						bestRow = i
						bestCol = j
					}
					alpha = max(alpha, bestScore)
				} else {
					if score < bestScore {
						bestScore = score
						bestRow = i
						bestCol = j
					}
					beta = min(beta, bestScore)
				}

				// Alpha-beta pruning
				if beta <= alpha {
					break
				}
			}
		}
	}

	// Return the best score along with the row and column of the best move
	return bestScore, bestRow, bestCol
}

// Helper function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
