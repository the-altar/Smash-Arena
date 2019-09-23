package server

import (
	"database/sql"
	"fmt"
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
		m := new(gameStart)
		if err := c.Bind(m); err != nil {

			fmt.Println("Here I am")
			return err
		}
		fmt.Println(m)
		return c.JSON(http.StatusOK, m)
	})

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	server.GET("/characters", getCharactersHandler)

	server.Logger.Fatal(server.Start(":" + port))
}
