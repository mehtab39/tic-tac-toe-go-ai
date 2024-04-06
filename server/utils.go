package server

import (
	"math"
	"os"
)

func GetAIUserId() string {
	return os.Getenv("USER_ID")
}

func isAITurn(gameStateUpdate *GameStateUpdate) bool {

	return gameStateUpdate.GameState.Players[gameStateUpdate.GameState.CurrentPlayer] == GetAIUserId()
}
func GetAIMove(game [3][3]string) (i, j int) {
	nextMoves := getEmptyCells(game)
	max_score := -1000000.0
	var max_move int
	for _, ij := range nextMoves {
		ij, score := minimax(game, ij, 0, true)
		if score > max_score {
			max_move = ij
			max_score = score
		}
	}

	i, j = unCellValue(max_move)
	return
}

func getEmptyCells(game [3][3]string) (cells []int) {
	for i, row := range game {
		for j, cell := range row {
			if cell == "" {
				cells = append(cells, cellValue(i, j))
			}
		}
	}
	return
}

func minimax(state [3][3]string, ij int, depth float64, maximise bool) (max_move int, max_score float64) {

	max_move = ij
	nextState := state
	i, j := unCellValue(ij)
	nextState[i][j] = playerToken(!maximise)

	winner := evaluate(nextState)
	switch winner {
	case playerToken(!maximise):
		max_score = 1.0
		return
	case playerToken(maximise):
		max_score = -1.0
		return
	case "":
		nextMoves := getEmptyCells(nextState)
		if len(nextMoves) == 0 {
			max_score = 0
			return
		}
		for _, xy := range nextMoves {
			_, score := minimax(nextState, xy, depth+1, !maximise)
			max_score = max_score + score*-1
		}
	}
	max_score = max_score * 0.5

	return
}

func unCellValue(ij int) (i, j int) {
	j = int(math.Mod(float64(ij), 10)) - 1
	i = (ij-j)/10 - 1
	return
}
func cellValue(i, j int) (ij int) {
	ij = (i+1)*10 + (j + 1)
	return
}

func playerToken(b bool) (c string) {
	if b {
		return "X"
	} else {
		return "O"
	}
}

func evaluate(game [3][3]string) (winner string) {
	wins := [][3]int{
		{11, 12, 13},
		{21, 22, 23},
		{31, 32, 33},
		{11, 21, 31},
		{12, 22, 32},
		{13, 23, 33},
		{11, 22, 33},
		{13, 22, 31}}

	var xs, ys []int
	for i, row := range game {
		for j, cell := range row {
			switch cell {
			case playerToken(true):
				xs = append(xs, cellValue(i, j))
			case playerToken(false):
				ys = append(ys, cellValue(i, j))
			}
		}
	}

	for _, winSet := range wins {
		if subset(winSet, xs) {
			winner = playerToken(true)
			return
		} else if subset(winSet, ys) {
			winner = playerToken(false)
			return
		} else {
			winner = ""
		}
	}

	return
}

func subset(first [3]int, second []int) bool {
	set := make(map[int]int)
	for _, value := range second {
		set[value] += 1
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}
