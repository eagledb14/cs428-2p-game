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
	currentPlayer := 1
	var row, col int

	for {
		move, quit := validateMsg(lobby)
		if quit {
			return
		}

		if move.Reset {
			board := types.NewFourInARowBoard()
			SendUpdate(lobby, board, currentPlayer, currentPlayer, true, false)
			continue
		}

		row, col = move.To.X, move.To.Y

		if isMoveValid(board, row, col, currentPlayer, move.Player) {
			board.Set(row, col, currentPlayer)

			if isGameOver(board, row, col, currentPlayer) {
				nextPlayer := ToggleRandomPlayer(2)
				SendUpdate(lobby, board, currentPlayer, nextPlayer, true, true)
				currentPlayer = nextPlayer
				continue
			}

			if isBoardFull(board) {
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
