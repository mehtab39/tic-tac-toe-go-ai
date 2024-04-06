package server

type GameStateUpdate struct {
	Type      string `json:"type"`
	GameState struct {
		ID            string       `json:"id"`
		Board         [3][3]string `json:"board"`
		CreatorID     string       `json:"creatorId"`
		Players       []string     `json:"players"`
		CurrentPlayer int          `json:"currentPlayer"`
		Status        string       `json:"status"`
		Winner        string       `json:"winner"`
	} `json:"gameState"`
}
