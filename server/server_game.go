package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// StartGameHandler will handle whenever a client wants to start a new game
func startGameHandler(c echo.Context) error {

	r := &startGameReq{}
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, 0)
	}

	team := buildTeam(r)
	gameRoom := buildGameRoom(team, r.UserID)
	arenas[r.UserID] = gameRoom // map the team that was just created so we can find it later
	return c.JSON(http.StatusOK, 1)
}

// GameRoomHandle will deal with everything else after the initial request is successful
func arenaHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	id := c.Param("id")

	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, 0)
	}
	g := gameHub{ws: ws, available: true, send: make(chan int), game: arenas[id]}

	go serveSocket(g)
	go listenSocket(g)

	return c.JSON(http.StatusOK, 1)
}
