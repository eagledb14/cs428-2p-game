package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "gopkg.in/olahol/melody.v1"
  "github.com/eagledb14/cs428-2p-game/types"
)

var router = gin.Default()

var socket = melody.New()

var lobbies = types.NewLobbyList()

func main() {
  fmt.Println("Hello World")

  router.GET("/:game", func(c *gin.Context) {
    id := c.Query("id")
    if id == "" {
      createLobby(c)
    } else {
      connectToLobby(c)
    }
  })

  l := types.NewLobby("checkers")
  fmt.Println(l)

  router.Run(":8080")
}
