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

	//set up blacks side,
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

	//set up red side
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
	return NewBoard(6, 7)
}

func (b *Board) Set(x int, y int, value int) error {
	index := x*b.column + y
	if index >= len(b.Board) || y >= b.row {
		return errors.New("Out of bounds Index")
	}

	b.Board[index] = value
	return nil
}

func (b *Board) Get(x int, y int) (int, error) {
	index := x*b.column + y
	if index >= len(b.Board) || y >= b.row {
		return 0, errors.New("Out of bounds Index")
	}

	return b.Board[index], nil
}

// func (b Board) GetBoard() []int {
//   return b.board
// }
