package types

import (
  "errors"
)

type Board struct {
  board []int
  row int
  column int
}

func NewBoard(row int, column int) Board {
  return Board {
    board: make([]int, row * column),
    row: row,
    column: column,
  }
}

func (b *Board) Set(x int, y int, value int) error {
  index := x * b.column + y
  if index >= len(b.board) || y >= b.row {
    return errors.New("Out of bounds Index")
  }

  b.board[index] = value
  return nil
}

func (b *Board) Get(x int, y int) (int, error) {
  index := x * b.column + y
  if index >= len(b.board) || y >= b.row {
    return 0, errors.New("Out of bounds Index")
  }

  return b.board[index], nil
}

// func (b Board) GetBoard() []int {
//   return b.board
// }


