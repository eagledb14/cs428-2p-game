package games

import (
  "github.com/eagledb14/cs428-2p-game/types"
  "encoding/json"
)

func validateMsg(lobby *types.Lobby) (types.Move, bool) {
  var move types.Move
  for {
    select {
    case <- lobby.Quit:
      return types.Move{}, true

    case msg := <- lobby.Chan:
      json.Unmarshal([]byte(msg), &move)
      return move, false

    default:
    }
  }
}


func SendUpdate(lobby *types.Lobby, board types.Board, currentPlayer int, validMove bool) {
  update := types.NewBoardUpdate(validMove, currentPlayer, board)
  json_update, _ := json.Marshal(update)
  for _, player := range lobby.Players {
    player.Write([]byte(json_update))
  }
}
