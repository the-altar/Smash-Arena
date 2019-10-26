package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Signup registers a new user
func Signup(g *gin.Context) {
	username := g.PostForm("username")
	password := g.PostForm("password")
	password, _ = hashPassword(password)

	if err := CreateUser(username, password); err != nil {
		g.Redirect(http.StatusMovedPermanently, "/")
	} else {
		g.Redirect(http.StatusMovedPermanently, "/signup")
	}

}

// Signin logs an user in
func Signin(g *gin.Context) {
	u := g.PostForm("username")
	p := g.PostForm("password")
	user, _ := OneUserByName(u)

	if checkPasswordHash(p, user.password) {
		g.Redirect(http.StatusMovedPermanently, "/")
	} else {
		g.Redirect(http.StatusMovedPermanently, "/login")
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
