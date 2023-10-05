package types

import (
	"gopkg.in/olahol/melody.v1"
)

type Lobby struct {
  GameType string
  Players []*melody.Session
  Chan chan string
  Quit chan struct{}
}

func NewLobby(gameType string) Lobby {
  return Lobby {
    GameType: gameType,
    Players: []*melody.Session{},
    Chan: make(chan string),
    Quit: make(chan struct{}),
  }
}
