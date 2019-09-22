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
func InitServer(port string, database *sql.DB) {
	db = database

	server.Static("/", "static")
	server.POST("/game", func(c echo.Context) error {
		m := echo.Map{}
		if err := c.Bind(&m); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, m)
	})

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	server.GET("/characters", getCharactersHandler)

	server.Logger.Fatal(server.Start(":" + port))
}
