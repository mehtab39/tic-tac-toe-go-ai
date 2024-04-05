package main

import "math"

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
