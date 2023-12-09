package games

import (
	"github.com/eagledb14/cs428-2p-game/types"
)

func Checkers(lobby *types.Lobby) {
	board := types.NewCheckersBoard()
	currentPlayer := 1
	var fromRow, toRow, fromCol, toCol int

	//these are used during multi-jumps
	requiredFromRow := -1
	requiredFromCol := -1
	jumpOnly := false

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

		if move.Pass {
			SendUpdate(lobby, board, currentPlayer, togglePlayer(currentPlayer), true, false)
			currentPlayer = togglePlayer(currentPlayer)
			continue
		}

		fromRow, toRow, fromCol, toCol = move.From.X, move.To.X, move.From.Y, move.To.Y

		//if move was a request for possible moves, send update with board showing possible moves
		//without ending the player's turn
		if move.GetMoves {
			_, moveBoard := getPossibleMoves(board, jumpOnly, fromCol, fromRow, requiredFromCol, requiredFromRow)
			SendUpdateSinglePlayer(lobby, moveBoard, currentPlayer, currentPlayer, currentPlayer, true, false)
			continue
		}

		if isCheckersMoveValid(board, jumpOnly, fromCol, toCol, fromRow, toRow, requiredFromCol, requiredFromRow, currentPlayer, move.Player) {
			pieceValue, _ := board.Get(fromCol, fromRow)

			//see if piece reached opposite end of board and should be promoted
			if shouldPieceBePromoted(pieceValue, toCol) {
				pieceValue += 2
			}

			//set origin space to be empty, update desination to have piece that moved
			board.Set(fromCol, fromRow, 0)
			board.Set(toCol, toRow, pieceValue)

			//check if move was a jump move
			jumpMove, jumpedCol, jumpedRow := getJumpedCoordinates(fromCol, toCol, fromRow, toRow)
			if jumpMove {
				//overwrite jumped piece with empty space
				board.Set(jumpedCol, jumpedRow, 0)

				//game can only end if a piece was eliminated
				if isCheckersGameOver(board) {
					SendUpdate(lobby, board, currentPlayer, togglePlayer(currentPlayer), true, true)
					continue
				}

				//check if the piece can jump again
				canJumpAgain, _ := getPossibleMoves(board, true, toCol, toRow, requiredFromCol, requiredFromRow)
				if canJumpAgain {
					//set variables to only allow jump moves from the current piece
					jumpOnly = true
					requiredFromRow = toRow
					requiredFromCol = toCol

					//update the player without ending their turn
					SendUpdateSinglePlayer(lobby, board, currentPlayer, currentPlayer, currentPlayer, true, false)
					continue
				}
			}

			//reset variables to allow any valid move from any piece
			jumpOnly = false
			requiredFromRow = -1
			requiredFromCol = -1

			//update both players and switch to next player's turn
			SendUpdate(lobby, board, currentPlayer, togglePlayer(currentPlayer), true, false)
			currentPlayer = togglePlayer(currentPlayer)
		} else {
			SendError(lobby, board, move, currentPlayer)
		}
	}

}

func getPossibleMoves(board types.Board, jumpOnly bool, currentCol, currentRow, requiredCol, requiredRow int) (bool, types.Board) {
	possibleMovesBoard := types.NewCheckersBoard()
	selectedPiece, _ := board.Get(currentCol, currentRow)
	moveFound := false

	//during multi-jumps, only the piece performing the jump is allowed to make moves
	if (requiredCol > -1 && requiredCol != currentCol) || (requiredRow > -1 && requiredRow != currentRow) {
		return moveFound, possibleMovesBoard
	}

	//normal pieces will only use the first check (red) or the second (black), but king pieces will use both
	if selectedPiece == 1 || selectedPiece > 2 {
		//get errors to check for out-of-bounds spaces before making other checks
		forwardLeft, forwardLeftError := board.Get(currentCol-1, currentRow-1)
		forwardRight, forwardRightError := board.Get(currentCol-1, currentRow+1)
		jumpLeft, jumpLeftError := board.Get(currentCol-2, currentRow-2)
		jumpRight, jumpRightError := board.Get(currentCol-2, currentRow+2)

		//if forward left is empty it can be moved to but not jumped
		if forwardLeftError == nil && forwardLeft == 0 && !jumpOnly {
			possibleMovesBoard.Set(currentCol-1, currentRow-1, 5)
			moveFound = true
		} else if jumpLeftError == nil && forwardLeft%2 != selectedPiece%2 && jumpLeft == 0 {
			//if it is occupied by an opponent piece and the space behind it is empty, it can be jumped
			possibleMovesBoard.Set(currentCol-2, currentRow-2, 5)
			moveFound = true
		}

		if forwardRightError == nil && forwardRight == 0 && !jumpOnly {
			possibleMovesBoard.Set(currentCol-1, currentRow+1, 5)
			moveFound = true
		} else if jumpRightError == nil && forwardRight%2 != selectedPiece%2 && jumpRight == 0 {
			possibleMovesBoard.Set(currentCol-2, currentRow+2, 5)
			moveFound = true
		}
	}

	if selectedPiece >= 2 {
		forwardLeft, forwardLeftError := board.Get(currentCol+1, currentRow-1)
		forwardRight, forwardRightError := board.Get(currentCol+1, currentRow+1)
		jumpLeft, jumpLeftError := board.Get(currentCol+2, currentRow-2)
		jumpRight, jumpRightError := board.Get(currentCol+2, currentRow+2)

		if forwardLeftError == nil && forwardLeft == 0 && !jumpOnly {
			possibleMovesBoard.Set(currentCol+1, currentRow-1, 5)
			moveFound = true
		} else if jumpLeftError == nil && forwardLeft%2 != selectedPiece%2 && jumpLeft == 0 {
			possibleMovesBoard.Set(currentCol+2, currentRow-2, 5)
			moveFound = true
		}

		if forwardRightError == nil && forwardRight == 0 && !jumpOnly {
			possibleMovesBoard.Set(currentCol+1, currentRow+1, 5)
		} else if jumpRightError == nil && forwardRight%2 != selectedPiece%2 && jumpRight == 0 {
			possibleMovesBoard.Set(currentCol+2, currentRow+2, 5)
			moveFound = true
		}
	}

	return moveFound, possibleMovesBoard
}

func isCheckersMoveValid(board types.Board, jumpOnly bool, fromCol, toCol, fromRow, toRow, requiredFromCol, requiredFromRow, currentPlayer, playerWhoMoved int) bool {
	if currentPlayer != playerWhoMoved {
		return false
	}

	selectedPiece, selectedPieceError := board.Get(fromCol, fromRow)

	//piece must be in-bounds and should either have same value as player (normal piece) or two greater (king piece)
	if selectedPieceError != nil || (selectedPiece != playerWhoMoved && selectedPiece != playerWhoMoved+2) {
		return false
	}

	//use the possible moves function to check if the attempted destination is valid
	_, possibleMovesBoard := getPossibleMoves(board, jumpOnly, fromCol, fromRow, requiredFromCol, requiredFromRow)
	destination, destinationError := possibleMovesBoard.Get(toCol, toRow)

	//if destination is valid, it will be marked on the board as a possible move by the number 5
	//if it is anything else, the move is invalid
	return destinationError == nil && destination == 5
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
boardCheck:
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
				break boardCheck
			}
		}
	}

	//if either team has no pieces on the board, the game is over
	return noRedPiecesLeft || noBlackPiecesLeft
}
