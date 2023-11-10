package sudoku

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
)

var Difficulties = []string{"easy", "medium", "hard"}
var gridValues = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}

const gridSize = 9
const boardSize = gridSize * gridSize

func PrintBoard(board [9][9]int) {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			fmt.Printf("%d ", board[row][col])
		}
		fmt.Println()
	}
}

func IsValidDifficulty(difficulty string) bool {
	for _, d := range Difficulties {
		if d == difficulty {
			return true
		}
	}
	return false
}

func IsValidBoard(board string) bool {
	if len(board) != boardSize {
		return false
	}

	for i := 0; i < boardSize; i++ {
		if board[i] != '.' {
			if board[i] < '1' || board[i] > '9' {
				return false
			}
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

func shuffle(a []int, rng *rand.Rand) []int {
	shuffled := make([]int, len(a))
	copy(shuffled, a)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

func fillDiagonal(board *[9][9]int, rng *rand.Rand) bool {
	for row := 0; row < gridSize; row++ {
		shift := row / 3 * 3
		for col := shift; col < 3+shift; col++ {
			shuffled := shuffle(gridValues[:], rng)

			for _, number := range shuffled {
				if isValidBox(board, row, col, number) {
					board[row][col] = number
					break
				}
			}
		}
	}
	return true
}

func checkIfFull(board *[9][9]int) bool {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if board[row][col] == 0 {
				return false
			}
		}
	}
	return true
}

func isValidBoard(board [9][9]int) bool {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if board[row][col] != 0 {
				number := board[row][col]
				board[row][col] = 0
				if !isValidPlacement(&board, row, col, number) {
					return false
				}
				board[row][col] = number
			}
		}
	}
	return true
}

func fillRemaining(board *[9][9]int, rng *rand.Rand) bool {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if board[row][col] == 0 {
				shuffled := shuffle(gridValues[:], rng)
				for _, number := range shuffled {
					if isValidPlacement(board, row, col, number) {
						board[row][col] = number
						if checkIfFull(board) {
							return true
						} else {
							if fillRemaining(board, rng) {
								return true
							}
						}
					}
				}
				return false
			}
		}
	}

	return true
}

func makeHoles(board *[9][9]int, difficulty string, rng *rand.Rand) {
	attempts := 0

	switch difficulty {
	case "easy":
		attempts = 3
	case "medium":
		attempts = 5
	case "hard":
		attempts = 7
	}

	for attempts > 0 {
		backup := 0
		row := 0
		col := 0

		for {
			row = rng.Intn(gridSize)
			col = rng.Intn(gridSize)

			if board[row][col] != 0 {
				backup = board[row][col]
				board[row][col] = 0
				break
			}
		}

		var copyBoard [9][9]int
		copy(copyBoard[:], board[:])
		counter := 0
		solver(&copyBoard, &counter)

		if counter != 1 {
			board[row][col] = backup
			attempts--
		}
	}
}

func GenerateBoard(difficulty string, seed string) (string, bool) {
	var board [9][9]int

	seedNumber, err := strconv.Atoi(seed)
	if err != nil {
		return "", false
	}

	rng := rand.New(rand.NewSource(int64(seedNumber)))

	ok := false
	for !ok {
		board = [9][9]int{}
		ok = fillDiagonal(&board, rng) && fillRemaining(&board, rng) && isValidBoard(board)
	}

	makeHoles(&board, difficulty, rng)
	PrintBoard(board)

	return convertBoardString(board), ok
}

func convertBoardArray(board string) [9][9]int {
	var boardArray [9][9]int

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if board[row*gridSize+col] == '.' {
				boardArray[row][col] = 0
			} else {
				val, err := strconv.Atoi(string(board[row*gridSize+col]))
				if err != nil {
					boardArray[row][col] = 0
				} else {
					boardArray[row][col] = val
				}
			}
		}
	}

	return boardArray
}

func convertBoardString(board [9][9]int) string {
	var boardString string

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if board[row][col] == 0 {
				boardString += "."
			} else {
				boardString += strconv.Itoa(board[row][col])
			}
		}
	}

	return boardString
}

func isValidRow(board *[9][9]int, row int, number int) bool {
	for i := 0; i < gridSize; i++ {
		if board[row][i] == number {
			return false
		}
	}
	return true
}

func isValidCol(board *[9][9]int, col int, number int) bool {
	for i := 0; i < gridSize; i++ {
		if board[i][col] == number {
			return false
		}
	}
	return true
}

func isValidBox(board *[9][9]int, row int, col int, number int) bool {
	startRow := row - row%3
	startCol := col - col%3

	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			if board[i][j] == number {
				return false
			}
		}
	}
	return true
}

func isValidPlacement(board *[9][9]int, row int, col int, number int) bool {
	return isValidBox(board, row, col, number) && isValidRow(board, row, number) && isValidCol(board, col, number)
}

func SolveBoard(board string) (string, int) {
	solution := convertBoardArray(board)
	counter := 0
	solver(&solution, &counter)
	return convertBoardString(solution), counter
}

func CounterToMsg(counter int) string {
	if counter == 0 {
		return "No solutions found"
	} else if counter == 1 {
		return "Unique solution found"
	} else {
		return "Multiple solutions found"
	}
}

func solver(board *[9][9]int, counter *int) {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if board[row][col] == 0 {
				for i := 1; i <= gridSize; i++ {
					if isValidPlacement(board, row, col, i) {
						board[row][col] = i
						solver(board, counter)
						board[row][col] = 0
					}
				}
				return
			}
		}
	}

	if checkIfFull(board) {
		*counter++
	}
}
