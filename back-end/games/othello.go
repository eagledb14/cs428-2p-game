package games

import (
	"fmt"

	"github.com/eagledb14/cs428-2p-game/types"
)

func Othello(lobby *types.Lobby) {
	board := types.NewOthelloBoard()
	currentPlayer := 1

	for {
		move, quit := validateMsg(lobby)
		if quit {
			return
		}

		if move.Reset {
			board = types.NewOthelloBoard()
			SendUpdate(lobby, board, currentPlayer, currentPlayer, true, false)
			continue
		}
		row, col := move.To.Y, move.To.X

		if isOthelloMoveValid(board, row, col, currentPlayer, move.Player) {
			updateOthelloBoard(&board, row, col, currentPlayer)

			if isOthelloOver(board, currentPlayer) {
				nextPlayer := ToggleRandomPlayer(2)
				SendUpdate(lobby, board, countOthelloWinner(board), nextPlayer, true, true)
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

func isOthelloMoveValid(board types.Board, row int, col int, currentPlayer int, playerWhoMoved int) bool {
	//checking correct player moved
	if currentPlayer != playerWhoMoved {
		return false
	}

	//check if piece is empty and move bounds
	if piece, err := board.Get(row, col); piece != 0 || err != nil {
		return false
	}

	otherPlayer := togglePlayer(currentPlayer)

	directions := []types.Point{
		{X: -1, Y: -1},
		{X: 0, Y: -1}, 
		{X: 1, Y: -1},
		{X: -1, Y: 0}, 
		{X: 1, Y: 0}, 
		{X: -1, Y: 1},
		{X: 0, Y: 1}, 
		{X: 1, Y: 1},
	}
	
	for _, dir := range directions {
		x, y := row+dir.X, col+dir.Y
		if piece, err := board.Get(x,y); err != nil || piece != otherPlayer {
			continue
		}

		// Move in the direction while we see the other player's pieces
		for piece, err := board.Get(x, y); err == nil && piece == otherPlayer; piece, err = board.Get(x, y) {
			x += dir.X
			y += dir.Y
		}

		// Check if we ended on a piece of the current player
		if piece, err := board.Get(x, y); err == nil && piece == currentPlayer {
			return true
		}
	}

	return false
}

func isOthelloOver(board types.Board, currentPlayer int) bool {

	//checks the whole board, if the place is empty, check if a move can be made, it not then return
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if piece, _ := board.Get(i, j); piece == 0 {
				if isOthelloMoveValid(board, i, j, currentPlayer, currentPlayer) {
					return false
				} 
			}
		}
	}

	return true
}

func countOthelloWinner(board types.Board) int {
	count1 := 0
	count2 := 0

	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if piece, _ := board.Get(i, j); piece == 1 {
				count1++
			} else if piece == 2 {
				count2++
			}
		}
	}

	if count1 > count2 {
		return 1
	} 
	return 2
}

func updateOthelloBoard(board *types.Board, row int, col int, currentPlayer int) {
	otherPlayer := togglePlayer(currentPlayer)
	    board.Set(row, col, currentPlayer)

		directions := []types.Point{
			{X: -1, Y: -1},
			{X: 0, Y: -1}, 
			{X: 1, Y: -1},
			{X: -1, Y: 0}, 
			{X: 1, Y: 0}, 
			{X: -1, Y: 1},
			{X: 0, Y: 1}, 
			{X: 1, Y: 1},
		}

	for _, dir := range directions {
		x, y := row+dir.X, col+dir.Y

		// Track if there are pieces to flip in this direction
		piecesToFlip := make([]types.Point, 0)

		// Move in the direction while we see the other player's pieces
		for piece, err := board.Get(x, y); err == nil && piece == otherPlayer; piece, err = board.Get(x, y) {
			piecesToFlip = append(piecesToFlip, types.Point{X: x, Y: y})
			x += dir.X
			y += dir.Y
		}

		// Check if we ended on a piece of the current player
		if piece, err := board.Get(x, y); err == nil && piece == currentPlayer {
			// Flip the pieces
			for _, p := range piecesToFlip {
				board.Set(p.X, p.Y, currentPlayer)
			}
		}
	}
}
