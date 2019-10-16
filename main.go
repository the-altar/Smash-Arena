package main

import (
	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/providers"
)

func main() {

	g := gin.New()
	g.Static("/public", "./public")
	g.StaticFile("/", "public/index.html")
	g.GET("ws/:id", func(g *gin.Context) {
		v := make(chan bool)
		go providers.Conn.Init(g, v)
		providers.Conn.PumpOut(g.Param("id"), <-v)
		return
	})

	g.Run()

}
