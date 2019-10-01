package server

import (
	"database/sql"
	"net/http"
	"smash/engine"

	"github.com/labstack/echo"
)

var (
	db     *sql.DB
	server *echo.Echo = echo.New()
)

// All built-in types we'll need in this package
type (
	clientMessageGame struct {
		Client string                   `json:"client"`
		Code   int                      `json:"code"`
		Data   map[int]engine.Character `json:"data"`
	}

	// startGameReq is the data from the initial request we get from the client to start a game
	startGameReq struct {
		UserID string   `json:"userID"`
		TeamID []string `json:"teamID"`
	}

	/* dbFeed contains all the info we get from the database. I'm using it as an struct to avoid passing
	a billion parameters to the constructor function*/
	dbFeed struct {
		charID           int
		charName         string
		skillID          int
		skillName        string
		skillDescription string
		effectName       string
		value            int
		duration         int
		tick             bool
	}
)

// InitServer starts the server
func InitServer(port string, dbase *sql.DB) {
	db = dbase
	server.Static("/", "static")

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	server.GET("/character", getCharactersHandler) // from server_character.go
	server.POST("/newgame", startGameHandler)      // from server_game.go
	server.GET("/arena", arenaHandler)             // from server_game.go
	server.HideBanner = true
	server.Logger.Fatal(server.Start(":" + port))
}
