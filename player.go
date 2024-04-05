package main

import "math"

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
