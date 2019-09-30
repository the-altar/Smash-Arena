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
func arenaHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, 0)
	}
	go socket(ws)
	return c.JSON(http.StatusOK, 1)
}
