package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/packages/providers"
	"golang.org/x/crypto/bcrypt"
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
		p := g.PostForm("password")
		hash, _ := hashPassword(p)
		fmt.Println(hash)
		fmt.Println(checkPasswordHash(p, hash))
		g.Redirect(http.StatusMovedPermanently, "/")

	})

	g.Run()

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
