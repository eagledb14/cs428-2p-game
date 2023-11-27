package games

import (
	"github.com/eagledb14/cs428-2p-game/types"
)

const (
	emptyCell = 0
	player1   = 1
	player2   = 2
	rows      = 6
	cols      = 7
)

func Fourinarow(lobby *types.Lobby) {
	board := types.NewFourInARowBoard()
	currentPlayer := player1

	for {
		move, quit := validateMsg(lobby)
		if quit {
			return
		}

		if move.Reset {
			board = types.NewFourInARowBoard()
			SendUpdate(lobby, board, currentPlayer, currentPlayer, true, false)
			continue
		}

		col := move.To.X

		board, piecePlaced := PlacePiece(board, col, move.Player, currentPlayer)

		if piecePlaced {
			if CheckWin(&board, currentPlayer) {
				nextPlayer := ToggleRandomPlayer(2)
				SendUpdate(lobby, board, currentPlayer, nextPlayer, true, true)
				currentPlayer = nextPlayer
				continue
			}
			if checkFull(board) {
				nextPlayer := ToggleRandomPlayer(2)
				SendUpdate(lobby, board, -1, nextPlayer, true, true)
				currentPlayer = nextPlayer
				continue
			}
			
			SendUpdate(lobby, board, currentPlayer, togglePlayer(currentPlayer), true, false)
			currentPlayer = togglePlayer(currentPlayer)
		} else {
			SendError(lobby, board, move, currentPlayer)
		}
	}
}

// PlacePiece places a piece on the board at the specified column for the given player.
// It returns the updated board and a boolean indicating whether the placement was successful.
func PlacePiece(board types.Board, column, player int, currentPlayer int) (types.Board, bool) {
	if player != currentPlayer {
		return board, false
	}

	// Iterate over rows to find the first empty cell in the specified column
	for row := rows - 1; row >= 0; row-- {
		cellValue, err := board.Get(row, column)
		if err != nil {
			return board, false
		}

		// If the cell is empty, place the piece and return the updated board
		if cellValue == emptyCell {
			board.Set(row, column, player)
			return board, true
		}
	}

	// If no empty cell is found in the column, return the original board and indicate failure
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

func checkFull(board types.Board) bool {
  for i := 0; i < cols; i++ {
    for j := 0; j < rows; j++ {
      if piece, _ := board.Get(j, i); piece == emptyCell {
	return false
      }
    }
  }

  return true
}
