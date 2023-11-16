package games

import (
	"testing"

	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/stretchr/testify/assert"
)

func TestTictactoe(t *testing.T) {

}

func TestCheckRow(t *testing.T) {
	// Create a sample 3x3 Tic Tac Toe board
	board := types.NewBoard(3, 3)

	// Fill a row with X symbols (row 1)
	for j := 0; j < 3; j++ {
		board.Set(1, j, 1)
	}

	// Check that row 1 is recognized as a win
	assert.True(t, checkRow(board, 1, 1))

	// Clear the board
	board = types.NewBoard(3, 3)

	// Fill a different row with O symbols (row 2)
	for j := 0; j < 3; j++ {
		board.Set(2, j, 2)
	}

	// Check that row 2 is recognized as a win
	assert.True(t, checkRow(board, 2, 2))
}

func TestIsMoveValid(t *testing.T) {
	// Create a sample 3x3 Tic Tac Toe board
	board := types.NewBoard(3, 3)

	// Define the current player and the player who just moved
	currentPlayer := 1
	playerWhoMoved := 1

	// Make some valid moves
	validMoves := [][2]int{
		{0, 0},
		{1, 1},
		{2, 2},
	}

	for _, move := range validMoves {
		row, col := move[0], move[1]
		assert.True(t, isTicTacToeMoveValid(board, row, col, currentPlayer, playerWhoMoved))
		board.Set(row, col, currentPlayer) // Simulate setting a piece
	}

	// Make some invalid moves
	invalidMoves := [][2]int{
		{0, 0},  // Already taken
		{1, 1},  // Already taken
		{-1, 0}, // Out of bounds
		{3, 0},  // Out of bounds
		{0, -1}, // Out of bounds
		{0, 3},  // Out of bounds
	}

	for _, move := range invalidMoves {
		row, col := move[0], move[1]
		assert.False(t, isTicTacToeMoveValid(board, row, col, currentPlayer, playerWhoMoved))
	}
}

func TestCheckDiagonals(t *testing.T) {
	// Create a sample 3x3 Tic Tac Toe board
	board := types.NewBoard(3, 3)

	// Fill the main diagonal with X symbols
	for i := 0; i < 3; i++ {
		board.Set(i, i, 1)
	}

	// Check that the main diagonal is recognized as a win
	assert.True(t, checkDiagonals(board, 1))

	// Clear the board
	board = types.NewBoard(3, 3)

	// Fill the other diagonal with O symbols
	for i := 0; i < 3; i++ {
		board.Set(i, 2-i, 2)
	}

	// Check that the other diagonal is recognized as a win
	assert.True(t, checkDiagonals(board, 2))
}

func TestIsBoardFull(t *testing.T) {
	// Create a sample 3x3 Tic Tac Toe board
	board := types.NewBoard(3, 3)

	// Initially, the board should not be full
	assert.False(t, isBoardFull(board))

	// Fill the entire board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board.Set(i, j, 1)
		}
	}

	// Now, the board should be full
	assert.True(t, isBoardFull(board))
}

func TestTogglePlayer(t *testing.T) {
	player_num := 1
	player_num = togglePlayer(player_num)
	if player_num != 2 {
		t.Errorf("Expected 2, got %d", player_num)
	}
	if togglePlayer(player_num) != 1 {
		t.Errorf("Expected 1, got %d", player_num)
	}
	if togglePlayer(420) != 1 {
		t.Errorf("Expected 1, got %d", player_num)
	}
}

func TestCheckColumn(t *testing.T) {
	// Create a sample 3x3 Tic Tac Toe board
	board := types.NewBoard(3, 3)

	// Fill a column with X symbols (column 0)
	for i := 0; i < 3; i++ {
		board.Set(i, 0, 1)
	}

	// Check that column 0 is recognized as a win
	assert.True(t, checkColumn(board, 0, 1))

	// Clear the board
	board = types.NewBoard(3, 3)

	// Fill a different column with O symbols (column 1)
	for i := 0; i < 3; i++ {
		board.Set(i, 1, 2)
	}

	// Check that column 1 is recognized as a win
	assert.True(t, checkColumn(board, 1, 2))
}
