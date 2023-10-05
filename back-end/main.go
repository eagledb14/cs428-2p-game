package main

import (
  "fmt"
  // "net/http"
  "github.com/gin-gonic/gin"
  "gopkg.in/olahol/melody.v1"
)

var router = gin.Default()

var socket = melody.New()

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


  router.Run(":8080")
}
