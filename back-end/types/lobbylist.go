package types

import (
  "sync"
)

type LobbyList struct {
  lobbies map[string]Lobby
  mutex sync.Mutex
}

func NewLobbyList() LobbyList {
  return LobbyList {
    lobbies: make(map[string]Lobby),
    mutex: sync.Mutex{},
  }
}

func (l *LobbyList) GetLobby(lobbyName string) (Lobby, bool) {
  l.mutex.Lock()
  defer l.mutex.Unlock()

  value, exists := l.lobbies[lobbyName]
  return value, exists
}

func (l *LobbyList) SetLobby(lobbyName string, lobby Lobby) {
  l.mutex.Lock()
  defer l.mutex.Unlock()

  l.lobbies[lobbyName] = lobby
}

func (l *LobbyList) RemoveLobby(lobbyName string) {
  l.mutex.Lock()
  defer l.mutex.Unlock()

  delete(l.lobbies, lobbyName)
}
