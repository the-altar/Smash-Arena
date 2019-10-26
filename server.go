package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/pkg/context/socket"
	"github.com/the-altar/Smash-Arena/pkg/context/user"
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

	g.GET("/signup", func(g *gin.Context) {
		g.HTML(http.StatusOK, "signup.gohtml", nil)
	})

	g.GET("/login", func(g *gin.Context) {
		g.HTML(http.StatusOK, "login.gohtml", nil)
	})

	g.GET("/ws/:id", socket.GameSocket)
	g.POST("/user/signin", user.Signin)
	g.POST("/user/signup", user.Signup)
	g.Run()
}
