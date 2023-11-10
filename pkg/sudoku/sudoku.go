package sudoku

import "regexp"

var Difficulties = []string{"easy", "medium", "hard"}

const gridSize = 9

func IsValidDifficulty(difficulty string) bool {
	for _, d := range Difficulties {
		if d == difficulty {
			return true
		}
	}
	return false
}

func IsValidBoard(board string) bool {
	if len(board) != gridSize*gridSize {
		return false
	}

	for i := 0; i < gridSize*gridSize; i++ {
		if board[i] < '1' || board[i] > '9' || board[i] == '.' {
			return false
		}
	}
	return true
}

func IsValidSeed(seed string) bool {
	pattern := "^[0-9]+$"

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return regex.MatchString(seed) && len(seed) < 20
}

func GenerateBoard(difficulty string, seed string) string {
	board := "123456789123456789123456789123456789123456789123456789123456789123456789123456789"
	return board
}

func SolveBoard(board string) string {
	solution := "123456789123456789123456789123456789123456789123456789123456789123456789123456789"
	return solution
}
