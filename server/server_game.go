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
	rManager.createRoom(r) // creates a room and adds it to the roomManager.Rooms map using the player's ID as key
	return c.JSON(http.StatusOK, 1)
}

// GameRoomHandle will deal with everything else after the initial request is successful
func arenaHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	id, chat := c.Param("id"), make(chan int)

	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, 0)
	}

	g := gameHub{ws: ws, available: true, send: make(chan int), game: rManager.Rooms[id]}

	rManager.addToPool(g)
	func() {
		if rManager.isFree() {
			return
		}
		rManager.makeBusy(true)
		go matchMaking()
	}()

	go serveSocket(g, chat)
	go listenSocket(g, id, chat)

	return c.JSON(http.StatusOK, 1)
}
