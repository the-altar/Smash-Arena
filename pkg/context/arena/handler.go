package arena

import (
	"github.com/gin-gonic/gin"
)

// GameSocket will handle our websocket connection
func GameSocket(g *gin.Context) {
	v := make(chan bool)
	go Conn.Begin(g, v)

	Conn.PumpOut(g.Param("id"), <-v)
}

// Arena serves our SPA
func Arena(g *gin.Context) {
	g.File("public/index.html")
}
