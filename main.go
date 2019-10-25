package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/providers"
)

func main() {

	g := gin.New()
	g.LoadHTMLGlob("templates/**/*")
	g.Static("/public", "./public")
	g.StaticFile("/arena", "public/index.html")

	g.GET("/", func(g *gin.Context) {
		g.HTML(http.StatusOK, "home.gohtml", gin.H{
			"title": "Main website",
		})
	})

	g.GET("/register", func(g *gin.Context) {
		g.HTML(http.StatusOK, "register.gohtml", nil)
	})

	g.GET("/ws/:id", func(g *gin.Context) {
		v := make(chan bool)
		go providers.Conn.Init(g, v)
		providers.Conn.PumpOut(g.Param("id"), <-v)
	})

	g.POST("/user/new", func(g *gin.Context) {

		fmt.Println(g.PostForm("username"))

		g.Redirect(http.StatusMovedPermanently, "/")

	})

	g.Run()

}
