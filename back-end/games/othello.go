package games

import (
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
		row, col := move.To.X, move.To.Y

		if isOthelloMoveValid(board, row, col, currentPlayer, move.Player) {
			updateOthelloBoard(&board, row, col, currentPlayer)

			if isOthelloOver(board) {
				nextPlayer := togglePlayer(2)
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

	//checking move bounds
	if row < 0 || row > 7 || col < 0 || col > 7 {
		return false
	}

	//check if piece is empty
	if piece, _ := board.Get(row, col); piece != 0 {
		return false
	}

	otherPlayer := togglePlayer(currentPlayer)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
				newX := row + dx
        newY := col + dy

        // Check bounds and then check for opponent's piece
				
			piece, _ := board.Get(newX, newY)
			if newX >= 0 && newX < 7 && newY >= 0 && newY < 7 && 
			piece == otherPlayer {
				return true
			}
		}
	}

	return false
}

func isOthelloOver(board types.Board) bool {

	//checks the whole board, if the place is empty, check if a move can be made, it not then return
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if piece, _ := board.Get(i, j); piece == 0 {
				if !isOthelloMoveValid(board, i, j, 1, 1) || !isOthelloMoveValid(board, i, j, 2, 2) {
					return true
				} 
				return false
			}
		}
	}

	return true
}

func countOthelloWinner(board types.Board) int {
	count1 := 0
	count2 := 0

	for i := 0; i <= 7; i++ {
		for j := 0; i <= 7; j++ {
			if piece, _ := board.Get(i, j); piece == 1 {
				count1++
			} else if piece == 1 {
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
	directions := []types.Point {
		{X: row-1, Y: col-1},
		{X: row, Y: col-1},  
		{X: row+1, Y: col-1},
		{X: row-1, Y: col},  
		{X: row+1, Y: col},  
		{X: row-1, Y: col+1},
		{X: row, Y: col+1},  
		{X: row+1, Y: col+1},
	}

	for _, dir := range(directions) {
		point := types.NewPoint(row + dir.X, col + dir.Y) 

		for piece, err := board.Get(point.X, point.Y); piece == otherPlayer && err != nil; {
			board.Set(point.X, point.Y, currentPlayer)
			point.AddPoint(dir)
		}
	}
}
