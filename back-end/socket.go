package main

import (
	"fmt"
	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/eagledb14/cs428-2p-game/games"
	"gopkg.in/olahol/melody.v1"
)


func handleSockets(m *melody.Melody) {
  m.HandleConnect(func(s *melody.Session) {
    handleConnect(s)
  })

  m.HandleMessage(func(s *melody.Session, msg []byte) {
    handleMessage(s, msg)
  })

  m.HandleDisconnect(func(s *melody.Session) {
    handleDisconnect(s)
  })
}

func handleConnect(s *melody.Session) {
  lobbyId := s.Request.URL.Query().Get("lobbyId")
  lobby, _ := lobbies.Get(lobbyId)

  //adding new connectiong to the lobby
  lobby.Players = append(lobby.Players, s)

  //checking if lobby size is the correct size
  if len(lobby.Players) != 2 {
    return
  }
  lobby.IsPlaying = true

  //sends each player their session id
  lobby.Shuffle()
  for i, player := range lobby.Players {
    player.Set("sessionId", i + 1)
    player.Write([]byte(fmt.Sprintf("%d", i + 1)))
  }

  switch lobby.GameType {
  case "tictactoe":
    go games.Tictactoe(lobby)
  case "checkers":
    go testGame(lobby)
  case "othello":
    go games.Othello(lobby)
  default:
    return
  }
}

func handleMessage(s *melody.Session, msg []byte) {
  //get lobby from lobbyList
  lobbyId := s.Request.URL.Query().Get("lobbyId")

  lobby, exists := lobbies.Get(lobbyId)
  if !exists {
    return
  }

  //sends message to the running game
  lobby.Chan <- string(msg)
}

func handleDisconnect(s *melody.Session) {
  lobbyId := s.Request.URL.Query().Get("lobbyId")

  lobby, exists := lobbies.Get(lobbyId)
  if !exists {
    return
  }

  //signaling quit to the running game
  lobby.Quit <- struct{}{}

  //closing the other players connections to the game
  for _, player := range lobby.Players {
    player.CloseWithMsg([]byte("Player Disconnected"))
  }

  //removing lobby from the available games
  lobbies.Remove(lobbyId)
}


func testGame(lobby *types.Lobby) {
  var msg string
  msg = <- lobby.Chan

  for _, player := range lobby.Players {
    player.Write([]byte(msg))
  }
}

