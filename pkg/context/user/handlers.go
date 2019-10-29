package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/the-altar/Smash-Arena/pkg/manager"
	"golang.org/x/crypto/bcrypt"
)

// Signup registers a new user
func Signup(g *gin.Context) {
	username := g.PostForm("username")
	password := g.PostForm("password")
	password, _ = hashPassword(password)

	if err := CreateUser(username, password); err != nil {
		user, err := OneUserByName(username)
		if err != nil {
			u, _ := uuid.NewUUID()
			g.SetCookie("sid", u.String(), 60*60*24, "/", "", false, true)
			manager.SetSession(u.String(), user.ID, user.Username)
		}
	}

	g.Redirect(http.StatusMovedPermanently, "/")
}

// Signin logs an user in
func Signin(g *gin.Context) {
	u := g.PostForm("username")
	p := g.PostForm("password")
	user, _ := OneUserByName(u)

	if checkPasswordHash(p, user.Password) {
		u, _ := uuid.NewUUID()
		g.SetCookie("sid", u.String(), 60*60*24, "/", "", false, true)
		manager.SetSession(u.String(), user.ID, user.Username)
	}

	g.Redirect(http.StatusMovedPermanently, "/")
}

// Signout removes all records from server
func Signout(g *gin.Context) {
	sid, _ := g.Cookie("sid")
	g.SetCookie("sid", "", -1, "/", "", false, true)
	manager.DestroySession(sid)

	g.Redirect(http.StatusMovedPermanently, "/")
}

// Self fetches the user from a session
func Self(g *gin.Context) {
	sid, err := g.Cookie("sid")
	if err != nil {
		return
	}
	session, _ := manager.GetSession(sid)
	g.JSON(http.StatusAccepted, session)
	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
