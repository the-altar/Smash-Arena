package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/pkg/context/account"
	"github.com/the-altar/Smash-Arena/pkg/context/admin"
	"github.com/the-altar/Smash-Arena/pkg/context/arena"
	"github.com/the-altar/Smash-Arena/pkg/context/home"
)

func main() {

	g := gin.New()
	g.LoadHTMLGlob("templates/**/*")
	g.GET("/", home.Home)

	g.GET("/admin", admin.Editor)
	g.POST("/admin/new/persona", arena.CreatePersona)
	g.GET("/admin/persona", arena.AllPersona)

	g.GET("/arena", arena.Arena)
	g.GET("/arena/ingame", arena.Arena)
	g.GET("/arena/ws/:id", arena.GameSocket)
	g.GET("/arena/api/persona", arena.AllPersona)
	g.GET("/arena/api/persona/skill/:id", arena.OneSkillSet)
	g.GET("/arena/api/account", account.Self)

	g.POST("/account/signin", account.Signin)
	g.POST("/account/signup", account.Signup)
	g.POST("/account/signout", account.Signout)

	g.Use(static.Serve("/", static.LocalFile("./public", true)))

	g.Run()

}
