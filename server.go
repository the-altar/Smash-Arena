package main

import (
	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/pkg/context/arena"
	"github.com/the-altar/Smash-Arena/pkg/context/home"
	"github.com/the-altar/Smash-Arena/pkg/context/user"
)

func main() {

	g := gin.New()
	g.LoadHTMLGlob("templates/**/*")
	g.Static("/public", "./public")

	g.GET("/", home.Home)

	g.GET("/arena", arena.Arena)
	g.GET("/arena/ws/:id", arena.GameSocket)
	g.GET("/arena/api/user", user.Self)

	g.POST("/user/signin", user.Signin)
	g.POST("/user/signup", user.Signup)
	g.POST("/user/signout", user.Signout)

	g.Run()
}
