package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/the-altar/Smash-Arena/providers"
)

func main() {

	Server := echo.New()
	Server.HideBanner = true

	Server.File("/", "public/index.html")

	Server.GET("ws/:id", func(c echo.Context) error {
		providers.Conn.Init(c)
		return nil
	})

	port := os.Getenv("PORT")

	Server.Logger.Fatal(Server.Start(":" + port))

}
