package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/the-altar/Smash-Arena/packages/providers"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func main() {

	providers.DB.Open()

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

	g.GET("/ws/:id", func(g *gin.Context) {
		v := make(chan bool)
		go providers.Conn.Init(g, v)
		providers.Conn.PumpOut(g.Param("id"), <-v)
	})

	g.POST("/user/new", func(g *gin.Context) {
		username := g.PostForm("username")
		password := g.PostForm("password")
		password, _ = hashPassword(password)

		providers.DB.CreateUser(username, password)

		g.Redirect(http.StatusMovedPermanently, "/")
	})

	g.POST("/user/login", func(g *gin.Context) {
		u := g.PostForm("username")
		p := g.PostForm("password")
		hash, _ := providers.DB.FindUserByName(u)

		if checkPasswordHash(p, hash) {
			g.Redirect(http.StatusMovedPermanently, "/")
		} else {
			g.Redirect(http.StatusMovedPermanently, "/login")
		}
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
