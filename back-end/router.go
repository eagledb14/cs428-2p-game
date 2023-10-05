package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  // "gopkg.in/olahol/melody.v1"
)

func createLobby(c *gin.Context) {
  c.String(http.StatusOK, "1235")
}

func connectToLobby(c *gin.Context) {
  c.String(http.StatusOK, "connected")
}


