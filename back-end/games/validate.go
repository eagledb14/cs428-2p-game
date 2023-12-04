package games

import (
	"encoding/json"
	"math/rand"

	"github.com/eagledb14/cs428-2p-game/types"
)

func validateMsg(lobby *types.Lobby) (types.Move, bool) {
	var move types.Move
	for {
		select {
		case <-lobby.Quit:
			return types.Move{}, true

		case msg := <-lobby.Chan:
			json.Unmarshal([]byte(msg), &move)
			return move, false

		default:
		}
	}
}

func SendUpdate(lobby *types.Lobby, board types.Board, currentPlayer, playerTurn int, validMove, finished bool) {
	update := types.NewBoardUpdate(validMove, currentPlayer, board, playerTurn)
	update.IsOver = finished
	json_update, _ := json.Marshal(update)
	for _, player := range lobby.Players {
		player.Write([]byte(json_update))
	}
}

func SendUpdateSinglePlayer(lobby *types.Lobby, board types.Board, currentPlayer, playerTurn, recipient int, validMove, finished bool) {
	update := types.NewBoardUpdate(validMove, currentPlayer, board, playerTurn)
	update.IsOver = finished
	json_update, _ := json.Marshal(update)
	lobby.Players[currentPlayer-1].Write([]byte(json_update))
}

func SendError(lobby *types.Lobby, board types.Board, move types.Move, playerTurn int) {
	update := types.NewBoardUpdate(false, move.Player, board, playerTurn)
	json_update, _ := json.Marshal(update)
	lobby.Players[move.Player-1].Write([]byte(json_update))
}

func ToggleRandomPlayer(numPlayers int) int {
	return rand.Intn(numPlayers) + 1
}
