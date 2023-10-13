package main

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/olahol/melody.v1"
  "github.com/eagledb14/cs428-2p-game/types"
)



var lobbies = types.NewLobbyList()

func main() {
  var router = gin.Default()
  var socket = melody.New()
  handleSockets(socket)

  // router.GET("/:game", func(c *gin.Context) {
  //   lobby := types.NewLobby("tictactoe")
  //   lobbies.Set("12345", &lobby)
  // })

  router.GET("/ws", func(c *gin.Context) {
    if lobbies.Len() < 1 {
      lobby := types.NewLobby("tictactoe")
      lobbies.Set("12345", &lobby)
    }
    socket.HandleRequest(c.Writer, c.Request)
  })
  router.Run(":8080")
}
