package games

import (
	"testing"

	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckers(t *testing.T) {

}

func TestValidMove(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//check that black piece is allowed to move forward
	value := isCheckersMoveValid(board, false, 2, 3, 3, 2, -1, -1, 2, 2)
	assert.True(t, value)
}

func TestGetPossibleMovesBlackPiece(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//check that black king piece is allowed to move backward
	_, moveBoard := getPossibleMoves(board, false, 2, 3, -1, -1)

	leftMoveValue, _ := moveBoard.Get(3, 2)
	rightMoveValue, _ := moveBoard.Get(3, 4)
	wrongMoveValue, _ := moveBoard.Get(4, 1)

	assert.True(t, leftMoveValue == 5)
	assert.True(t, rightMoveValue == 5)
	assert.True(t, wrongMoveValue == 0)
}

func TestGetPossibleMovesRedPiece(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//check that black king piece is allowed to move backward
	_, moveBoard := getPossibleMoves(board, false, 5, 2, -1, -1)

	leftMoveValue, _ := moveBoard.Get(4, 1)
	rightMoveValue, _ := moveBoard.Get(4, 3)
	wrongMoveValue1, _ := moveBoard.Get(3, 4)
	wrongMoveValue2, _ := moveBoard.Get(3, 6)

	assert.True(t, leftMoveValue == 5)
	assert.True(t, rightMoveValue == 5)
	assert.True(t, wrongMoveValue1 == 0)
	assert.True(t, wrongMoveValue2 == 0)
}

func TestValidKingMove(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//set empty space behind black piece and promote black piece to king
	board.Set(1, 2, 0)
	board.Set(2, 3, 4)

	//check that black king piece is allowed to move backward
	value := isCheckersMoveValid(board, false, 2, 1, 3, 2, -1, -1, 2, 2)
	assert.True(t, value)
}

func TestInvalidMove(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//set empty space behind black piece
	board.Set(1, 2, 0)

	//check that black piece is not allowed to move backward
	value := isCheckersMoveValid(board, false, 2, 1, 3, 2, -1, -1, 2, 2)
	assert.False(t, value)
}

func TestValidJump(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//move black piece forward
	board.Set(2, 3, 0)
	board.Set(3, 2, 2)

	//move red piece forward, in postion to jump black piece
	board.Set(5, 2, 0)
	board.Set(4, 1, 1)

	//check that red piece is allowed to jump black piece
	value := isCheckersMoveValid(board, false, 4, 2, 1, 3, -1, -1, 1, 1)
	assert.True(t, value)
}

func TestInvalidJump(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//move red piece forward
	board.Set(5, 0, 0)
	board.Set(4, 1, 1)

	//set up a red king piece where it could jump the red piece
	board.Set(3, 2, 3)

	//check that red king piece is not allowed to jump red piece
	value := isCheckersMoveValid(board, false, 3, 5, 2, 0, -1, -1, 1, 1)
	assert.False(t, value)
}

func TestGetJumpedPiece(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//move black piece forward
	board.Set(2, 3, 0)
	board.Set(3, 4, 2)

	//move red piece forward, in postion to be jumped by black piece
	board.Set(5, 2, 0)
	board.Set(4, 3, 1)

	//check that logic correctly identifies space to be jumped,
	//and that the coordinates are that of the red piece being jumped
	value, jumpedCol, jumpedRow := getJumpedCoordinates(3, 5, 4, 2)
	jumpedValue, _ := board.Get(jumpedCol, jumpedRow)

	assert.True(t, value)
	assert.True(t, jumpedCol == 4)
	assert.True(t, jumpedRow == 3)
	assert.True(t, jumpedValue == 1)
}

func TestPromotePieces(t *testing.T) {
	//try to promote a normal red piece in the correct place -- should return true
	value1 := shouldPieceBePromoted(1, 0)
	//try to promote a normal black piece in the correct place -- should return true
	value2 := shouldPieceBePromoted(2, 7)
	//try to promote a king red piece in the correct place -- should return false
	value3 := shouldPieceBePromoted(3, 0)
	//try to promote a normal black piece in the wrong place -- should return false
	value4 := shouldPieceBePromoted(2, 5)

	assert.True(t, value1)
	assert.True(t, value2)
	assert.False(t, value3)
	assert.False(t, value4)
}

func TestNoWin(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//check that win condition is not met
	value := isCheckersGameOver(board)
	assert.False(t, value)
}

func TestWinTeam1(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//replace all of team two's pieces with empty spaces
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pieceValue, _ := board.Get(i, j)
			if pieceValue == 2 {
				board.Set(i, j, 0)
			}
		}
	}

	//check that this triggers a win condition
	value := isCheckersGameOver(board)
	assert.True(t, value)
}

func TestWinTeam2(t *testing.T) {
	// Create a new checkers board
	board := types.NewCheckersBoard()

	//replace all of team one's pieces with empty spaces
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			pieceValue, _ := board.Get(i, j)
			if pieceValue == 1 {
				board.Set(i, j, 0)
			}
		}
	}

	//check that this triggers a win condition
	value := isCheckersGameOver(board)
	assert.True(t, value)
}
