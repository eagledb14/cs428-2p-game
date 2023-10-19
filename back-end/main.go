package main

import (
	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

var lobbies = types.NewLobbyList()

func main() {
	var router = gin.Default()
	var socket = melody.New()
	handleSockets(socket)
	handleRoutes(router, socket)
}
