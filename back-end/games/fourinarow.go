package games

import (
	"fmt"

	"github.com/eagledb14/cs428-2p-game/types"
)

const (
	emptyCell = 0
	player1   = 1
	player2   = 2
	rows      = 6
	cols      = 7
)

func PrintBoard(b *types.Board) {
	for _, row := range b.Board {
		fmt.Println(row)
	}
}

func Print(b *types.Board) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			value, err := b.Get(i, j)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Print(value, " ")
		}
		fmt.Println()
	}
}

func Fourinarow(lobby *types.Lobby) {
	board := types.NewFourInARowBoard()
	player := player1
	winner := 0

	fmt.Println("Welcome to Four-In-A-Row!")
	PrintBoard(&board)

	for winner == 0 {
		fmt.Printf("Player %d, choose a column (0-6): ", player)
		var col int
		fmt.Scanln(&col)

		if col < 0 || col >= cols {
			fmt.Println("Invalid column. Choose a column between 0 and 6.")
			continue
		}
		board, piecePlaced := PlacePiece(board, col, player)

		if piecePlaced {
			PrintBoard(&board)
			if CheckWin(&board, player) {
				winner = player
			} else if player == player1 {
				player = player2
			} else {
				player = player1
			}
		} else {
			fmt.Println("Column is full. Choose another column.")
		}
	}

	if winner > 0 {
		fmt.Printf("Player %d wins!\n", winner)
	} else {
		fmt.Println("It's a draw. The game is over.")
	}
}

func PlacePiece(board types.Board, column, player int) (types.Board, bool) {
	rows, err := board.Get(0, 0)
	if err != nil {
		return board, false
	}

	if column < 0 || column >= rows {
		return board, false
	}

	for row := rows - 1; row >= 0; row-- {
		value, err := board.Get(row, column)
		if err != nil {
			return board, false
		}

		if value == emptyCell {
			board.Set(row, column, player)
			return board, true
		}
	}

	return board, false
}

func CheckWin(b *types.Board, player int) bool {
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell, _ := b.Get(row, col)
			if cell == player {
				// Check horizontally
				if col+3 < cols {
					value1, err1 := b.Get(row, col+1)
					value2, err2 := b.Get(row, col+2)
					value3, err3 := b.Get(row, col+3)

					if err1 == nil && err2 == nil && err3 == nil &&
						value1 == player && value2 == player && value3 == player {
						return true
					}
				}

				// Check vertically
				if row+3 < rows {
					value1, err1 := b.Get(row+1, col)
					value2, err2 := b.Get(row+2, col)
					value3, err3 := b.Get(row+3, col)

					if err1 == nil && err2 == nil && err3 == nil &&
						value1 == player && value2 == player && value3 == player {
						return true
					}
				}

				// Check diagonally (up-right)
				if row-3 >= 0 && col+3 < cols {
					value1, err1 := b.Get(row-1, col+1)
					value2, err2 := b.Get(row-2, col+2)
					value3, err3 := b.Get(row-3, col+3)

					if err1 == nil && err2 == nil && err3 == nil &&
						value1 == player && value2 == player && value3 == player {
						return true
					}
				}

				// Check diagonally (up-left)
				if row-3 >= 0 && col-3 >= 0 {
					value1, err1 := b.Get(row-1, col-1)
					value2, err2 := b.Get(row-2, col-2)
					value3, err3 := b.Get(row-3, col-3)

					if err1 == nil && err2 == nil && err3 == nil &&
						value1 == player && value2 == player && value3 == player {
						return true
					}
				}

			}
		}
	}
	return false
}
