package games

import (
	"fmt"
)

const (
	emptyCell = 0
	player1   = 1
	player2   = 2
	rows      = 6
	cols      = 7
)

type Board struct {
	board [][]int
}

func NewBoard() *Board {
	board := make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
	}
	return &Board{board}
}

func (b *Board) Print() {
	for _, row := range b.board {
		fmt.Println(row)
	}
}

func (b *Board) PlacePiece(column, player int) bool {
	if column < 0 || column >= cols || b.board[0][column] != emptyCell {
		return false
	}

	for row := rows - 1; row >= 0; row-- {
		if b.board[row][column] == emptyCell {
			b.board[row][column] = player
			return true
		}
	}

	return false
}

func (b *Board) CheckWin(player int) bool {
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if b.board[row][col] == player {
				// Check horizontally
				if col+3 < cols &&
					b.board[row][col+1] == player &&
					b.board[row][col+2] == player &&
					b.board[row][col+3] == player {
					return true
				}

				// Check vertically
				if row+3 < rows &&
					b.board[row+1][col] == player &&
					b.board[row+2][col] == player &&
					b.board[row+3][col] == player {
					return true
				}

				// Check diagonally (up-right)
				if row-3 >= 0 && col+3 < cols &&
					b.board[row-1][col+1] == player &&
					b.board[row-2][col+2] == player &&
					b.board[row-3][col+3] == player {
					return true
				}

				// Check diagonally (up-left)
				if row-3 >= 0 && col-3 >= 0 &&
					b.board[row-1][col-1] == player &&
					b.board[row-2][col-2] == player &&
					b.board[row-3][col-3] == player {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	board := NewBoard()
	player := player1
	winner := 0

	fmt.Println("Welcome to Four-In-A-Row!")
	board.Print()

	for winner == 0 {
		fmt.Printf("Player %d, choose a column (0-6): ", player)
		var col int
		fmt.Scanln(&col)

		if col < 0 || col >= cols {
			fmt.Println("Invalid column. Choose a column between 0 and 6.")
			continue
		}

		if board.PlacePiece(col, player) {
			board.Print()
			if board.CheckWin(player) {
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
