package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// InitServer boots our server
func InitServer() {
	e := echo.New()
	e.Static("/", "static")
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
