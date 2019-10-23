package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/providers"
)

func main() {

	g := gin.New()
	g.LoadHTMLGlob("./templates/*")
	g.Static("/public", "./public")
	g.StaticFile("/arena", "public/index.html")

	g.GET("/", func(g *gin.Context) {
		g.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Main website",
		})
	})

	g.GET("/register", func(g *gin.Context) {
		g.HTML(http.StatusOK, "register.html", nil)
	})

	g.GET("ws/:id", func(g *gin.Context) {
		v := make(chan bool)
		go providers.Conn.Init(g, v)
		providers.Conn.PumpOut(g.Param("id"), <-v)
	})

	g.Run()

}
