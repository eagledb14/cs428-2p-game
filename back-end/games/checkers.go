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

func isCheckersMoveValid(board types.Board, fromCol, toCol, fromRow, toRow, currentPlayer, playerWhoMoved int) bool {
	if currentPlayer != playerWhoMoved {
		return false
	}

	//each row and col parameter must correspond to one of the 8 rows and cols
	if fromRow > 7 || fromRow < 0 || toRow > 7 || toRow < 0 || fromCol > 7 || fromCol < 0 || toCol > 7 || toCol < 0 {
		return false
	}

	//will return -1 if coordinates are out of bounds
	selectedPiece, _ := board.Get(fromCol, fromRow)
	destination, _ := board.Get(toCol, toRow)

	//selected space must be in-bounds and have a piece, destination must be in-bounds and empty
	if selectedPiece < 1 || destination != 0 {
		return false
	}

	//piece should either have same value as player (normal piece) or two greater (king piece)
	if selectedPiece != playerWhoMoved && selectedPiece != playerWhoMoved+2 {
		return false
	}

	if selectedPiece == 1 {
		//move check for normal red pieces
		if toCol >= fromCol || toCol < fromCol-2 {
			//cannot move backwards, stay in current column, or move more than two columns forwards
			return false
		} else if toCol == fromCol-1 && !(toRow == fromRow+1 || toRow == fromRow-1) {
			//destination being one column forward from origin indicates normal move,
			//in which case the destination row must be 1 less or greater than origin
			return false
		} else if toCol == fromCol-2 {
			//execution will only reach here if destination is two columns forward from origin, indicating a jump move
			var jumpedSpace int

			//when jumping, the destination row must be two greater or less than the origin row
			if toRow == fromRow+2 {
				jumpedSpace, _ = board.Get(fromCol-1, fromRow+1)
			} else if toRow == fromRow-2 {
				jumpedSpace, _ = board.Get(fromCol-1, fromRow-1)
			} else {
				return false
			}

			//space being jumped must contain either a normal piece or king piece from black team
			if !(jumpedSpace == 2 || jumpedSpace == 4) {
				return false
			}
		}
	} else if selectedPiece == 2 {
		//move check for normal black pieces
		if toCol <= fromCol || toCol > fromCol+2 {
			//cannot move backwards, stay in current column, or move more than two columns forwards
			return false
		} else if toCol == fromCol+1 && !(toRow == fromRow+1 || toRow == fromRow-1) {
			//destination being one column forward from origin indicates normal move,
			//in which case the destination row must be 1 less or greater than origin
			return false
		} else if toCol == fromCol+2 {
			//execution will only reach here if destination is two columns forward from origin, indicating a jump move
			var jumpedSpace int

			//when jumping, the destination row must be two greater or less than the origin row
			if toRow == fromRow+2 {
				jumpedSpace, _ = board.Get(fromCol+1, fromRow+1)
			} else if toRow == fromRow-2 {
				jumpedSpace, _ = board.Get(fromCol+1, fromRow-1)
			} else {
				return false
			}

			//space being jumped must contain either a normal piece or king piece from red team
			if !(jumpedSpace == 1 || jumpedSpace == 3) {
				return false
			}
		}

	} else {
		//move check for any king piece
		if toCol > fromCol+2 || toCol < fromCol-2 {
			//cannot stay in current column, or move more than two columns forwards or backwards
			return false
		} else if (toCol == fromCol-1 || toCol == fromCol+1) && !(toRow == fromRow+1 || toRow == fromRow-1) {
			//destination being one column forward or back from origin indicates normal move,
			//in which case the destination row must be 1 less or greater than origin
			return false
		} else if toCol == fromCol-2 || toCol == fromCol+2 {
			//execution will only reach here if destination is two columns forward or back from origin, indicating a jump move
			var jumpedSpace int

			//there are four possible spaces a king could be jumping to at any time,
			//so the destination must have the coordinates of one of those four
			if toCol == fromCol+2 && toRow == fromRow+2 {
				jumpedSpace, _ = board.Get(fromCol+1, fromRow+1)
			} else if toCol == fromCol+2 && toRow == fromRow-2 {
				jumpedSpace, _ = board.Get(fromCol+1, fromRow-1)
			} else if toCol == fromCol-2 && toRow == fromRow+2 {
				jumpedSpace, _ = board.Get(fromCol-1, fromRow+1)
			} else if toCol == fromCol-2 && toRow == fromRow-2 {
				jumpedSpace, _ = board.Get(fromCol-1, fromRow-1)
			} else {
				return false
			}

			//space being jumped must contain either a normal piece or king piece from oopposing team
			if jumpedSpace == 0 || jumpedSpace%2 == playerWhoMoved%2 {
				return false
			}
		}
	}

	return true
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
				break
			}
		}
		if !noRedPiecesLeft && !noBlackPiecesLeft {
			break
		}
	}

	//if either team has no pieces on the board, the game is over
	return noRedPiecesLeft || noBlackPiecesLeft
}
