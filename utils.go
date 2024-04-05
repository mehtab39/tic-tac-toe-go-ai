package main

import (
	"strings"
)

func ExtractGameIDFromURL(path string) (string, bool) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return "", false
	}
	return parts[2], true
}

func IsMyTurn(gameStateUpdate *GameStateUpdate) bool {
	return gameStateUpdate.GameState.Players[gameStateUpdate.GameState.CurrentPlayer] == AIGoUserID
}
