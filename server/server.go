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

// InitServer boots our server
func InitServer(port string, dbase *sql.DB) {
	db = dbase
	server.Static("/", "static")

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	server.GET("/characters", getCharactersHandler)
	server.POST("/newgame", startGameHandler)

	server.Logger.Fatal(server.Start(":" + port))
}
