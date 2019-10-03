package server

import (
	"database/sql"
	"net/http"
	"smash/engine"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	db       *sql.DB
	server   *echo.Echo = echo.New()
	upgrader            = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	arenas   = make(map[string]*engine.GameRoom)
	gamePool = make([]gameHub, 0)
	mutex    = &sync.Mutex{}
)

// All built-in types we'll need in this package
type (
	// charClient is what we send back to the client when they request information about a character
	charClient struct {
		ID      int    `json:"ID"`
		Name    string `json:"Name"`
		Profile string `json:"Profile"`
	}
	// clientMessageGame is what we expect to get from the client once they're in game
	clientMessageGame struct {
		Client    string                   `json:"client"`
		Code      int                      `json:"code"`
		GameState map[int]engine.Character `json:"gameState"`
	}

	// startGameReq is the data from the initial request we get from the client to start a game
	startGameReq struct {
		UserID string   `json:"userID"`
		TeamID []string `json:"teamID"`
	}

	// A game hub for each game
	gameHub struct {
		available bool
		ws        *websocket.Conn
		send      chan int
		game      *engine.GameRoom
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
	server.Use(middleware.CORS())
	server.Static("/", "static")
	server.File("/", "static/index.html")
	server.GET("/character", getCharactersHandler) // from server_character.go
	server.POST("/newgame", startGameHandler)      // from server_game.go
	server.GET("/arena/:id", arenaHandler)         // from server_game.go
	server.HideBanner = true
	server.Logger.Fatal(server.Start(":" + port))

	go func() {
		for {
			matchmake()
			time.Sleep(10 * time.Second)
		}
	}()
}
