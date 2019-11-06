package main

import (
	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/pkg/context/account"
	"github.com/the-altar/Smash-Arena/pkg/context/arena"
	"github.com/the-altar/Smash-Arena/pkg/context/home"
)

func main() {

	g := gin.New()
	g.LoadHTMLGlob("templates/**/*")
	g.Static("/public", "./public")

	g.GET("/", home.Home)

	g.GET("/arena", arena.Arena)
	g.GET("/arena/ws/:id", arena.GameSocket)
	g.GET("/arena/api/persona", arena.GetAllPersona)
	g.GET("/arena/api/account", account.Self)

	g.POST("/account/signin", account.Signin)
	g.POST("/account/signup", account.Signup)
	g.POST("/account/signout", account.Signout)

	g.Run()
}
