package server

import (
	"net/http"
	"smash/engine"

	"github.com/labstack/echo"
)

func startGameHandler(c echo.Context) error {
	var game map[int]engine.Character

	r := &startGameReq{}

	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, 0)
	}
	game = buildTeam(r) // from ./server_helpers.go
	return c.JSON(http.StatusOK, game)
}
