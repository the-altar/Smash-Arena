package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/pkg/context/account"
	"github.com/the-altar/Smash-Arena/pkg/manager"
)

// Home page handler
func Home(g *gin.Context) {
	cookie, _ := g.Cookie("sid")
	if session, ok := manager.GetSession(cookie); ok {
		u, _ := account.OneAccountByID(session.ID)
		g.HTML(http.StatusOK, "home.html", gin.H{
			"user": u,
		})
	} else {
		g.HTML(http.StatusOK, "home.html", nil)
	}
}
