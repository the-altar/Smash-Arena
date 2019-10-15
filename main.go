package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/the-altar/Smash-Arena/providers"
)

func main() {

	Server := echo.New()
	Server.Logger.SetLevel(log.OFF)
	Server.HideBanner = true

	Server.File("/", "public/index.html")

	Server.GET("ws/:id", func(c echo.Context) error {
		return providers.Conn.Init(c)
	})

	Server.Logger.Fatal(Server.Start(":8080"))

}
