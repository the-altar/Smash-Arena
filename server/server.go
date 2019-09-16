package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// InitServer boots our server
func InitServer(port string) {
	e := echo.New()

	e.Static("/", "static")
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	/*var engine engine.GameRoom
	team := [3]string{"N", "R", "T"}
	engine.Begin("Joao", team)*/
	e.Logger.Fatal(e.Start(":" + port))
}
