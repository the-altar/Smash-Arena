package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/pkg/context/arena"
	"github.com/the-altar/Smash-Arena/pkg/context/user"
	"github.com/the-altar/Smash-Arena/pkg/manager"
)

func main() {

	g := gin.New()
	g.LoadHTMLGlob("templates/**/*")
	g.Static("/public", "./public")

	g.GET("/", func(g *gin.Context) {
		cookie, _ := g.Cookie("sid")
		if session, ok := manager.GetSession(cookie); ok {
			u, _ := user.OneUserByID(session.ID)
			g.HTML(http.StatusOK, "home.html", gin.H{
				"user": u,
			})
		} else {
			g.HTML(http.StatusOK, "home.html", nil)
		}
	})

	g.GET("/signup", func(g *gin.Context) {
		g.HTML(http.StatusOK, "signup.html", nil)
	})

	g.GET("/login", func(g *gin.Context) {
		g.HTML(http.StatusOK, "signin.html", nil)
	})

	g.GET("/arena", arena.Arena)
	g.GET("/arena/ws/:id", arena.GameSocket)
	g.GET("/arena/api/user", user.Self)
	g.POST("/user/signin", user.Signin)
	g.POST("/user/signup", user.Signup)
	g.POST("/user/signout", user.Signout)

	g.Run()
}
