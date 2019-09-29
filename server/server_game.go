package server

import (
	"net/http"
	"smash/engine"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	arenas = make(map[string]*engine.GameRoom)
)

// All built-in types we'll need in this file
type (

	// startGameReq is the data from the initial request we get from the client to start a game
	startGameReq struct {
		UserID string   `json:"userID"`
		TeamID []string `json:"teamID"`
	}

	/* dbFeed contains all the info we get from the database. I'm using it as an struct to avoid passing
	a billion parameters to the constructor function*/
	dbFeed struct {
		charID           int
		charName         string
		skillID          int
		skillName        string
		skillDescription string
		effectName       string
		value            int
		duration         int
		tick             bool
	}
)

// StartGameHandler will handle whenever a client wants to start a new game
func startGameHandler(c echo.Context) error {

	r := &startGameReq{}
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, 0)
	}

	servant := buildTeam(r)
	groom := buildGameRoom(servant, r.UserID)
	arenas[r.UserID] = groom

	return c.JSON(http.StatusOK, 1)
}

// GameRoomHandle will deal with everything else after the initial request is successful
func gameRoomHandle(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, 0)
	}
	go socket(ws)
	return c.JSON(http.StatusOK, 1)
}
