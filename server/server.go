package server

import (
	"database/sql"
	"fmt"
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
	rManager = roomManager{}
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

	roomManager struct {
		Rooms     map[string]*engine.GameRoom
		GamePool  []gameHub
		freeRooms bool
	}
)

// InitServer starts the server
func InitServer(port string, dbase *sql.DB) {
	db = dbase
	rManager.begin()

	server.Use(middleware.CORS())
	server.Static("/", "static")
	server.File("/", "static/index.html")
	server.GET("/character", getCharactersHandler) // from server_character.go
	server.POST("/newgame", startGameHandler)      // from server_game.go
	server.GET("/arena/:id", arenaHandler)         // from server_game.go
	server.HideBanner = true
	server.Logger.Fatal(server.Start(":" + port))
}

func matchMaking() {
	for {
		fmt.Println("Matchmaking...")
		matchmake()
		time.Sleep(10 * time.Second)

		if rManager.poolSize() < 1 {
			fmt.Println("Killing process...")
			rManager.makeBusy(false)
			return
		}
	}
}

func (r *roomManager) begin() {
	r.Rooms = make(map[string]*engine.GameRoom)
	r.GamePool = make([]gameHub, 0)
	r.freeRooms = false
}

func (r *roomManager) createRoom(req *startGameReq) {
	room := buildGameRoom(req)
	r.Rooms[req.UserID] = room
}

func (r *roomManager) addToPool(g gameHub) {
	mutex.Lock()
	r.GamePool = append(r.GamePool, g)
	mutex.Unlock()
}

func (r *roomManager) poolSize() int {
	mutex.Lock()
	size := len(r.GamePool)
	mutex.Unlock()
	return size
}

func (r *roomManager) makeBusy(f bool) {
	mutex.Lock()
	r.freeRooms = f
	mutex.Unlock()
}

func (r roomManager) isFree() bool {
	mutex.Lock()
	f := r.freeRooms
	mutex.Unlock()
	return f
}

func (r *roomManager) poolPop() (int, gameHub) {
	s := r.poolSize()
	if s > 0 {
		mutex.Lock()
		x := r.GamePool[s-1]
		r.GamePool = r.GamePool[:s-1]
		mutex.Unlock()

		return 1, x
	}

	return 0, gameHub{}
}

func (gh *gameHub) joinEnemy(t map[int]engine.Character, pID string) {
	gh.game.SetOpponent(pID)
	gh.game.AddEnemies(t)
	gh.game.Full <- true
}
