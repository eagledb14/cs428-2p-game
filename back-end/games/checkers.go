package games

import (
	"github.com/eagledb14/cs428-2p-game/types"
)

func Checkers(lobby *types.Lobby) {
	board := types.NewCheckersBoard()
	currentPlayer := 1
	var fromRow, toRow, fromCol, toCol int

	for {
		move, quit := validateMsg(lobby)
		if quit {
			return
		}

		if move.Reset {
			board = types.NewCheckersBoard()
			SendUpdate(lobby, board, currentPlayer, currentPlayer, true, false)
			continue
		}

		fromRow, toRow, fromCol, toCol = move.From.X, move.To.X, move.From.Y, move.To.Y

		//if move was a request for possible moves, send update with board showing possible moves
		//without ending the player's turn
		if move.GetMoves {
			SendUpdate(lobby, getPossibleMoves(board, fromCol, fromRow), currentPlayer, currentPlayer, true, false)
			continue
		}

		if isCheckersMoveValid(board, fromCol, toCol, fromRow, toRow, currentPlayer, move.Player) {
			pieceValue, _ := board.Get(fromCol, fromRow)

			//see if piece reached opposite end of board and should be promoted
			if shouldPieceBePromoted(pieceValue, toCol) {
				pieceValue += 2
			}

			//check if move was a jump move, set jumped space to empty if it was
			jumpMove, jumpedCol, jumpedRow := getJumpedCoordinates(fromRow, toRow, fromCol, toCol)
			if jumpMove {
				board.Set(jumpedCol, jumpedRow, 0)
			}

			//set origin space to be empty, update desination to have piece that moved
			board.Set(fromCol, fromRow, 0)
			board.Set(toCol, toRow, pieceValue)

			if isCheckersGameOver(board) {
				//set player to random
				SendUpdate(lobby, board, currentPlayer, togglePlayer(currentPlayer), true, true)
				continue
			}

			SendUpdate(lobby, board, currentPlayer, togglePlayer(currentPlayer), true, false)
			currentPlayer = togglePlayer(currentPlayer)
		} else {
			SendError(lobby, board, move, currentPlayer)
		}
	}

}

func getPossibleMoves(board types.Board, currentCol, currentRow int) types.Board {
	possibleMovesBoard := types.NewCheckersBoard()
	selectedPiece, _ := board.Get(currentCol, currentRow)

	//normal pieces will only use the first check (red) or the second (black), but king pieces will use both
	if (selectedPiece == 1 || selectedPiece > 2) && currentCol != 0 {
		forwardLeft, _ := board.Get(currentCol-1, currentRow-1)
		forwardRight, _ := board.Get(currentCol-1, currentRow+1)
		jumpLeft, _ := board.Get(currentCol-2, currentRow-2)
		jumpRight, _ := board.Get(currentCol-2, currentRow+2)

		//if forward left is empty and in bounds, it can be moved to but not jumped
		if currentRow != 0 && forwardLeft == 0 {
			possibleMovesBoard.Set(currentCol-1, currentRow-1, 5)
		} else if forwardLeft == 2 && jumpLeft == 0 && currentCol > 1 && currentRow > 1 {
			//check for king
			//if it is occupied by an opponent piece and the space behind it is empty and in-bounds, it can be jumped
			possibleMovesBoard.Set(currentCol-2, currentRow-2, 5)
		}

		if currentRow != 7 && forwardRight == 0 {
			possibleMovesBoard.Set(currentCol-1, currentRow+1, 5)
		} else if forwardRight == 2 && jumpRight == 0 && currentCol > 1 && currentRow < 6 {
			possibleMovesBoard.Set(currentCol-2, currentRow+2, 5)
		}
	}

	if selectedPiece >= 2 && currentCol != 7 {
		forwardLeft, _ := board.Get(currentCol+1, currentRow-1)
		forwardRight, _ := board.Get(currentCol+1, currentRow+1)
		jumpLeft, _ := board.Get(currentCol+2, currentRow-2)
		jumpRight, _ := board.Get(currentCol+2, currentRow+2)

		if currentRow != 0 && forwardLeft == 0 {
			possibleMovesBoard.Set(currentCol+1, currentRow-1, 5)
		} else if forwardLeft == 2 && jumpLeft == 0 && currentCol < 6 && currentRow > 1 {
			possibleMovesBoard.Set(currentCol+2, currentRow-2, 5)
		}

		if currentRow != 7 && forwardRight == 0 {
			possibleMovesBoard.Set(currentCol+1, currentRow+1, 5)
		} else if forwardRight == 2 && jumpRight == 0 && currentCol < 6 && currentRow < 6 {
			possibleMovesBoard.Set(currentCol+2, currentRow+2, 5)
		}
	}

	return possibleMovesBoard
}

func isCheckersMoveValid(board types.Board, fromCol, toCol, fromRow, toRow, currentPlayer, playerWhoMoved int) bool {
	if currentPlayer != playerWhoMoved {
		return false
	}

	//each row and col parameter must be in-bounds
	//use get errors for bound checking
	if fromRow > 7 || fromRow < 0 || toRow > 7 || toRow < 0 || fromCol > 7 || fromCol < 0 || toCol > 7 || toCol < 0 {
		return false
	}

	selectedPiece, _ := board.Get(fromCol, fromRow)

	//piece should either have same value as player (normal piece) or two greater (king piece)
	if selectedPiece != playerWhoMoved && selectedPiece != playerWhoMoved+2 {
		return false
	}

	//use the possible moves function to check if the attempted destination is valid
	possibleMovesBoard := getPossibleMoves(board, fromCol, fromRow)
	destination, _ := possibleMovesBoard.Get(toCol, toRow)

	//if destination is valid, it will be marked on the board as a possible move by the number 5
	//if it is anything else, the move is invalid
	return destination == 5
}

func shouldPieceBePromoted(pieceValue, destCol int) bool {
	if pieceValue > 2 {
		//pieces is already a king
		return false
	} else if pieceValue == 1 {
		//black piece, make king if it reached bottom of the board
		return destCol == 0
	} else {
		//red piece, make king if it reached top of the board
		return destCol == 7
	}
}

func getJumpedCoordinates(fromCol, toCol, fromRow, toRow int) (bool, int, int) {
	//return bool indicating if move is a jump and coordinates of space being jumped
	if toCol == fromCol+2 && toRow == fromRow+2 {
		return true, fromCol + 1, fromRow + 1
	} else if toCol == fromCol+2 && toRow == fromRow-2 {
		return true, fromCol + 1, fromRow - 1
	} else if toCol == fromCol-2 && toRow == fromRow+2 {
		return true, fromCol - 1, fromRow + 1
	} else if toCol == fromCol-2 && toRow == fromRow-2 {
		return true, fromCol - 1, fromRow - 1
	} else {
		return false, -1, -1
	}
}

func isCheckersGameOver(board types.Board) bool {
	noRedPiecesLeft := true
	noBlackPiecesLeft := true

	//look in each space for a piece
checkBoard:
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pieceValue, _ := board.Get(i, j)
			if pieceValue == 1 || pieceValue == 3 {
				noRedPiecesLeft = false
			} else if pieceValue == 2 || pieceValue == 4 {
				noBlackPiecesLeft = false
			}

			//no need to keep searching if a piece from each team has been found
			if !noRedPiecesLeft && !noBlackPiecesLeft {
				break checkBoard
			}
		}
	}

	//if either team has no pieces on the board, the game is over
	return noRedPiecesLeft || noBlackPiecesLeft
}
