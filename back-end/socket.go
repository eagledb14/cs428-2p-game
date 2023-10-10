package main

import (
	"fmt"
	"gopkg.in/olahol/melody.v1"
  "github.com/eagledb14/cs428-2p-game/types"
)


func handleSockets(m *melody.Melody) {
  m.HandleConnect(func(s *melody.Session) {
    handleConnect(s)
  })
}

func handleConnect(s *melody.Session) {
  lobbyId := s.Request.URL.Query().Get("id")
  lobby, _ := lobbies.Get(lobbyId)
  
  s.Set("lobbyId", lobbyId)
  s.Set("gameType", lobby.GameType)

  //adding new connectiong to the lobby
  lobby.Players = append(lobby.Players, s)

  //checking if lobby size is the correct size
  if len(lobby.Players) != 2 {
    return
  }

  //sends each player their session id
  lobby.Shuffle()
  for i, player := range lobby.Players {
    player.Set("sessionId", i)
    player.Write([]byte(fmt.Sprintf("%d", i)))
  }

  switch lobby.GameType {
  case "tictactoe":
    go testGame(lobby)
  case "checkers":
    go testGame(lobby)
  default:
    return
  }
}


func testGame(lobby *types.Lobby) {
  for i, player := range lobby.Players {
    player.Write([]byte(fmt.Sprintf("connected to new lobby: %d", i)))
  }
}

