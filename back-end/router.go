package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

// contains logic for all routes
func handleRoutes(router *gin.Engine, socket *melody.Melody) {
	//initializes a new lobby for a game or connects to an existing one if an id query is given
	router.GET("/:game", func(c *gin.Context) {
		id := c.Query("lobbyId")
		game := c.Param("game")
		board, gameExists := types.NewGame(game)

		//make sure the provided game parameter corresponds with a valid game type
		if gameExists {
			//if no id was given, create a new lobby
			if id == "" {
				lobby := types.NewLobby(game)
				lobbyID := createLobbyID(c)

				//loops until an unused lobby id is found
				for _, ok := lobbies.Get(lobbyID); ok; _, ok = lobbies.Get(lobbyID) {
					lobbyID = createLobbyID(c)
				}

				//add lobby with id to list and send id to frontend
				lobbies.Set(lobbyID, &lobby)
				c.String(http.StatusOK, lobbyID)
			} else {
				//if an id was given, try to find corresponding lobby in list
				_, lobbyFound := lobbies.Get(id)

				//if lobby exists, create new board and send to frontend
				if lobbyFound {
					json_update, _ := json.Marshal(board)
					c.String(http.StatusOK, string(json_update))
				} else {
					//if lobby doesn't exist, notify frontend
					c.String(http.StatusOK, "Lobby not found")
				}
			}
		} else {
			//notify frontend if game parameter was invalid
			c.String(http.StatusOK, game+" is not a supported game")
		}
	})

	router.GET("/ws", func(c *gin.Context) {
		lobbyID := c.Query("lobbyId")

		//lobby id must be provided
		if lobbyID != "" {
			lobby, lobbyFound := lobbies.Get(lobbyID)
			//make sure lobby exists
			if lobbyFound {
				//if game is not already in progress, upgrade http request to web socket
				if !lobby.IsPlaying {
					socket.HandleRequest(c.Writer, c.Request)
				} else {
					//if game is in progress, don't let client join
					c.String(http.StatusOK, "Game already in progress")
				}
			} else {
				c.String(http.StatusOK, "Lobby not found")
			}
		} else {
			c.String(http.StatusOK, "No lobby ID")
		}
	})
	router.Run(":8080")
}

func createLobbyID(c *gin.Context) string {
	lobbyID := strconv.Itoa(rand.Intn(999999))
	return lobbyID
}
