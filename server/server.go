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
		Client    string    `json:"client"`
		Code      int       `json:"code"`
		GameState gameState `json:"gameState"`
	}

	gameState struct {
		Opponent string           `json:"opponent"`
		Foes     map[string]_char `json:"foes"`
		Friends  map[string]_char `json:"friends"`
	}

	_char struct {
		ID     int                `json:"id"`
		Health int                `json:"health"`
		Skills map[string]_skills `json:"skills"`
	}

	_skills struct {
		ID int `json:"id"`
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
		ongoing   chan bool
		Game      *engine.GameRoom `json:"gameState"`
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
		Rooms    map[string]*engine.GameRoom
		GamePool struct {
			games   []*gameHub
			players map[string]bool
		}
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
	r.GamePool.games, r.GamePool.players = make([]*gameHub, 0), make(map[string]bool)
	r.freeRooms = false
}

func (r *roomManager) deleteRoom(id string) {
	delete(r.Rooms, id)
}

func (r *roomManager) createRoom(req *startGameReq) {
	room := buildGameRoom(req)
	r.Rooms[req.UserID] = &room
}

func (r *roomManager) removeFromPool(id string) {
	mutex.Lock()
	_, ok := r.GamePool.players[id]
	if ok {
		for i, n := range r.GamePool.games {
			if n.Game.Player == id {
				delete(r.GamePool.players, n.Game.Player)
				r.GamePool.games = append(r.GamePool.games[:i], r.GamePool.games[i+1:]...)
				mutex.Unlock()
				return
			}
		}
	} else {
		mutex.Unlock()
		return
	}
}

func (r *roomManager) poolAppend(g *gameHub) {
	mutex.Lock()
	r.GamePool.games = append(r.GamePool.games, g)
	r.GamePool.players[g.Game.GetPlayer()] = true
	mutex.Unlock()
}

func (r *roomManager) poolSize() int {
	mutex.Lock()
	size := len(r.GamePool.games)
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

func (r *roomManager) poolPop() (int, *gameHub) {
	s := r.poolSize()
	if s > 0 {
		mutex.Lock()

		x := r.GamePool.games[s-1]
		delete(r.GamePool.players, x.Game.GetPlayer())
		r.GamePool.games = r.GamePool.games[:s-1]

		fmt.Println("popped, players left: ", r.GamePool.players)
		mutex.Unlock()

		return 1, x
	}
	return 0, &gameHub{}
}

func (gh *gameHub) joinEnemy(t map[string]engine.Character, pID string) {
	gh.Game.SetOpponent(pID)
	gh.Game.AddEnemies(t)
	gh.ongoing <- true
}

func (cm *clientMessageGame) writeGameState(g *engine.GameRoom) {
	cm.Client = g.Player
	cm.Code = 1
	cm.GameState = gameState{g.Opponent, make(map[string]_char), make(map[string]_char)}

	for key, char := range g.GetTeam() {
		cm.GameState.Friends[key] = _char{char.ID, char.Health, make(map[string]_skills)}
		for sKey, skill := range char.Skills {
			cm.GameState.Friends[key].Skills[sKey] = _skills{skill.ID}
		}
	}

	for key, char := range g.Enemies {
		cm.GameState.Foes[key] = _char{char.Health, char.ID, make(map[string]_skills)}
	}
}
