package admin

import (
	"github.com/gin-gonic/gin"
)

// Editor handler
func Editor(g *gin.Context) {
	g.File("./public/admin/index.html")
	return
}
