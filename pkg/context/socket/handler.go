package socket

import (
	"github.com/gin-gonic/gin"
)

// GameSocket will handle our websocket connection
func GameSocket(g *gin.Context) {
	v := make(chan bool)
	go Conn.Begin(g, v)

	Conn.PumpOut(g.Param("id"), <-v)
}
