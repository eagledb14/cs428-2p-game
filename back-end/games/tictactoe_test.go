package games

import (
	"testing"

	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/stretchr/testify/assert"
)

func TestTictactoe(t *testing.T) {

}

func TestCheckRow(t *testing.T) {

}

func TestIsMoveValid(t *testing.T) {

}

func TestCheckDiagonals(t *testing.T) {

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

}
