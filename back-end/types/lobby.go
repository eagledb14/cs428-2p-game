package types

import (
	"gopkg.in/olahol/melody.v1"
	"math/rand"
	"time"
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
    Chan: make(chan string, 10),
    Quit: make(chan struct{}),
  }
}


func (l *Lobby) Shuffle() {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  n := len(l.Players)

  for i := n - 1; i > 0; i-- {
    j := r.Intn(i + 1)
    l.Players[i], l.Players[j] = l.Players[j], l.Players[i]
  }
}
