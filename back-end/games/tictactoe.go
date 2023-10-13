package games

import (
	"fmt"
	"encoding/json"
	"github.com/eagledb14/cs428-2p-game/types"
)

func Tictactoe(lobby *types.Lobby) {
	board := types.NewTicTacToeBoard()
	currentPlayer := 1
	var row, col int

	for {
		move, quit := validateMsg(lobby)
		if quit {
			return
		}
		row, col = move.To.X, move.To.Y

		if isMoveValid(board, row, col, currentPlayer, move.Player) {
			board.Set(row, col, currentPlayer)

			if isGameOver(board, row, col, currentPlayer) {
				SendUpdate(lobby, board, currentPlayer, true, true)
				continue
			}

			if isBoardFull(board) {
				fmt.Println("It's a tie!")
				continue
			}
			SendUpdate(lobby, board, currentPlayer, true, false)
			currentPlayer = togglePlayer(currentPlayer)
		} else {
			update := types.NewBoardUpdate(false, move.Player, board)
			json_update, _ := json.Marshal(update)
			lobby.Players[move.Player].Write([]byte(json_update))
		}
	}
}

func isMoveValid(board types.Board, row, col int, currentPlayer int, playerWhoMoved int) bool {
	if currentPlayer-1 != playerWhoMoved {
		return false
	}
	cell, _ := board.Get(row, col)
	return cell == 0
}

func isGameOver(board types.Board, row, col, currentPlayer int) bool {
	return checkRow(board, row, currentPlayer) || checkColumn(board, col, currentPlayer) || checkDiagonals(board, currentPlayer)
}

func checkRow(board types.Board, row, currentPlayer int) bool {
	for i := 0; i < 3; i++ {
		cell, _ := board.Get(row, i)
		if cell != currentPlayer {
			return false
		}
	}
	return true
}

func checkColumn(board types.Board, col, currentPlayer int) bool {
	for i := 0; i < 3; i++ {
		cell, _ := board.Get(i, col)
		if cell != currentPlayer {
			return false
		}
	}
	return true
}

func checkDiagonals(board types.Board, currentPlayer int) bool {
	// Check the main diagonal (top-left to bottom-right)
	for i := 0; i < 3; i++ {
		cell, _ := board.Get(i, i)
		if cell != currentPlayer {
			break
		}
		if i == 2 {
			return true
		}
	}

	// Check the other diagonal (top-right to bottom-left)
	for i := 0; i < 3; i++ {
		cell, _ := board.Get(i, 2-i)
		if cell != currentPlayer {
			break
		}
		if i == 2 {
			return true
		}
	}

	return false
}

func isBoardFull(board types.Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cell, _ := board.Get(i, j)
			if cell == 0 {
				return false
			}
		}
	}
	return true
}

func togglePlayer(currentPlayer int) int {
	if currentPlayer == 1 {
		return 2
	}
	return 1
}
