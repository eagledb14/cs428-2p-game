package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/eagledb14/cs428-2p-game/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gopkg.in/olahol/melody.v1"
)

// contains logic for all routes
func handleRoutes(router *gin.Engine, socket *melody.Melody) {
	router.Use(cors.Default())

	//initializes a new lobby for a game or connects to an existing one if an id query is given
	router.GET("/:game", func(c *gin.Context) {
		id := c.Query("lobbyId")
		gameType := c.Param("game")
		board, gameExists := types.NewGame(gameType)

		if !gameExists {
			//notify frontend if gameType parameter was invalid
			c.String(http.StatusOK, gameType + " is not a supported game")
			return
		}

		//if no id was given, create a new lobby
		if id == "" {
			handleNewLobby(c, gameType)
		} else {
			handleLobbyConnect(c, gameType, id, board)
		}
	})

	//upgrades a clients http connection to a web socket
	router.GET("/ws", func(c *gin.Context) {
		lobbyID := c.Query("lobbyId")

		//lobby id is required
		if lobbyID == "" {
			c.String(http.StatusOK, "No lobby ID")
			return
		}

		lobby, lobbyFound := lobbies.Get(lobbyID)

		//return error msg is lobby doesn't exist or game is already in session
		if !lobbyFound {
			c.String(http.StatusOK, "Lobby not found")
			return
		} else if lobby.IsPlaying {
			c.String(http.StatusOK, "Game already in progress")
			return
		}

		//upgrade http connection to web socket
		socket.HandleRequest(c.Writer, c.Request)
	})

	//returns the type of game being hosted in a lobby
	router.GET("/lobby", func(c *gin.Context) {
		lobbyID := c.Query("lobbyId")

		//lobby id is required
		if lobbyID == "" {
			c.String(http.StatusOK, "No lobby ID")
			return
		}

		lobby, lobbyFound := lobbies.Get(lobbyID)

		//return error msg is lobby doesn't exist
		if !lobbyFound {
			c.String(http.StatusOK, "Lobby not found")
			return
		}

		//return lobby game type
		c.String(http.StatusOK, lobby.GameType)

	})

	router.Run(":8080")
}

func handleNewLobby(c *gin.Context, gameType string) {
	lobby := types.NewLobby(gameType)
	lobbyID := createLobbyID(c)

	//loops until an unused lobby id is found
	for _, ok := lobbies.Get(lobbyID); ok; _, ok = lobbies.Get(lobbyID) {
		lobbyID = createLobbyID(c)
	}

	//add lobby with id to list and send id to frontend
	lobbies.Set(lobbyID, &lobby)
	c.String(http.StatusOK, lobbyID)
}

func handleLobbyConnect(c *gin.Context, gameType string, lobbyId string, board types.Board) {
	//if an id was given, try to find corresponding lobby in list
	lobby, lobbyFound := lobbies.Get(lobbyId)

	//send error msg if lobby doesn't exist or if game types are mistmatched
	if !lobbyFound {
		c.String(http.StatusOK, "Lobby not found")
		return
	} 

	if lobby.IsPlaying {
		c.String(http.StatusOK, "Lobby already in play")
		return
	} 

	if lobby.GameType != gameType {
		c.String(http.StatusOK, "Wrong game type")
		return
	} 

	//create new board and send to frontend
	json_update, _ := json.Marshal(board)
	c.String(http.StatusOK, string(json_update))
}

func createLobbyID(c *gin.Context) string {
	lobbyID := strconv.Itoa(rand.Intn(999999))
	return lobbyID
}
