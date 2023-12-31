package types

import (
	"errors"
)

type Board struct {
	Board  []int `json:"board"`
	row    int
	column int
}

func NewGame(gameType string) (Board, bool) {
	switch gameType {
	case "tictactoe":
		return NewTicTacToeBoard(), true
	case "checkers":
		return NewCheckersBoard(), true
	case "fourinarow":
		return NewFourInARowBoard(), true
	case "othello":
		return NewOthelloBoard(), true
	}

	return NewBoard(0, 0), false
}

func NewBoard(row int, column int) Board {
	return Board{
		Board:  make([]int, row*column),
		row:    row,
		column: column,
	}
}

func NewCheckersBoard() Board {
	checkersBoard := NewBoard(8, 8)

	//spaces on the board that cannot have a piece on them are marked with -1
	for i := range checkersBoard.Board {
		checkersBoard.Board[i] = -1
	}

	//set black team pieces
	//1, 3, 5, 7
	checkersBoard.Set(0, 1, 2)
	checkersBoard.Set(0, 3, 2)
	checkersBoard.Set(0, 5, 2)
	checkersBoard.Set(0, 7, 2)
	//0, 2, 4, 6
	checkersBoard.Set(1, 0, 2)
	checkersBoard.Set(1, 2, 2)
	checkersBoard.Set(1, 4, 2)
	checkersBoard.Set(1, 6, 2)
	//1, 3, 5, 7
	checkersBoard.Set(2, 1, 2)
	checkersBoard.Set(2, 3, 2)
	checkersBoard.Set(2, 5, 2)
	checkersBoard.Set(2, 7, 2)

	//set blank spaces in middle of board
	//0, 2, 4, 6
	checkersBoard.Set(3, 0, 0)
	checkersBoard.Set(3, 2, 0)
	checkersBoard.Set(3, 4, 0)
	checkersBoard.Set(3, 6, 0)
	//1, 3, 5, 7
	checkersBoard.Set(4, 1, 0)
	checkersBoard.Set(4, 3, 0)
	checkersBoard.Set(4, 5, 0)
	checkersBoard.Set(4, 7, 0)

	//set red team pieces
	//0, 2, 4, 6
	checkersBoard.Set(5, 0, 1)
	checkersBoard.Set(5, 2, 1)
	checkersBoard.Set(5, 4, 1)
	checkersBoard.Set(5, 6, 1)
	//1, 3, 5, 7
	checkersBoard.Set(6, 1, 1)
	checkersBoard.Set(6, 3, 1)
	checkersBoard.Set(6, 5, 1)
	checkersBoard.Set(6, 7, 1)
	//0, 2, 4, 6
	checkersBoard.Set(7, 0, 1)
	checkersBoard.Set(7, 2, 1)
	checkersBoard.Set(7, 4, 1)
	checkersBoard.Set(7, 6, 1)

	return checkersBoard
}

func NewTicTacToeBoard() Board {
	return NewBoard(3, 3)
}

func NewFourInARowBoard() Board {
	return NewBoard(7, 6)
}

func NewOthelloBoard() Board {
	othelloBoard := NewBoard(8,8)

	//set black squares
	othelloBoard.Set(3, 3, 1)
	othelloBoard.Set(4, 4, 1)

	//set white squares
	othelloBoard.Set(3, 4, 2)
	othelloBoard.Set(4, 3, 2)

	return othelloBoard
}


func (b *Board) Set(x int, y int, value int) error {
	index := x*b.row + y
	if index >= len(b.Board) || y >= b.row || x >= b.column || x < 0 || y < 0 {
		return errors.New("Out of bounds Index")
	}

	b.Board[index] = value
	return nil
}

func (b *Board) Get(x int, y int) (int, error) {
	index := x*b.row + y
	if index >= len(b.Board) || y >= b.row || x >= b.column || x < 0 || y < 0 {
		return 0, errors.New("Out of bounds Index")
	}

	return b.Board[index], nil
}

// func (b Board) GetBoard() []int {
//   return b.board
// }
