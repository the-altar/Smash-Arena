package server

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

var (
	db     *sql.DB
	server *echo.Echo = echo.New()
)

// InitServer starts the server
func InitServer(port string, dbase *sql.DB) {
	db = dbase
	server.Static("/", "static")

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	server.GET("/character", getCharactersHandler) // from api_character.go
	server.GET("/newgame", startGameHandler)       // from api_game.go
	server.Logger.Fatal(server.Start(":" + port))
}
