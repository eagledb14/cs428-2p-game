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

  m.HandleMessage(func(s *melody.Session, msg []byte) {
    handleMessage(s, msg)
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

func handleMessage(s *melody.Session, msg []byte) {
  //get lobby from lobbyList
  lobbyId, _ := s.Get("lobbyId")
  lobby, _ := lobbies.Get(lobbyId.(string))

  //sends message to the running game
  lobby.Chan <- string(msg)
}


func testGame(lobby *types.Lobby) {
  var msg string
  msg = <- lobby.Chan

  for _, player := range lobby.Players {
    player.Write([]byte(msg))
  }
}

