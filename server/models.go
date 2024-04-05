package server

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

// Represents the state of the Tic-Tac-Toe game
type GameState struct {
	Board         [][]string // 3x3 board
	CurrentPlayer string     // Current player ('X' or 'O')
	Maximizing    bool       // Indicates if it's the maximizing player's turn
}
