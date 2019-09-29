package server

import (
	"fmt"
	"net/http"
	"smash/engine"
	"smash/gamedb"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type (
	// startGameReq is the data from the initial request we get from the client to start a game
	startGameReq struct {
		UserID string   `json:"UserID"`
		TeamID []string `json:"TeamID"`
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

// StartGameHandler will handle whenever a client wants to start a new game
func startGameHandler(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	go socket(ws)
	return c.JSON(http.StatusOK, 1)
}

// BuildTeam builds the player's team for the game
func buildTeam(r *startGameReq) map[int]engine.Character {
	charMap := make(map[int]engine.Character)
	teamQuery := gamedb.QueryCharData(r.TeamID)
	rows, err := db.Query(teamQuery)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		res := &dbFeed{}
		err = rows.Scan(&res.charID, &res.charName, &res.skillID, &res.skillName, &res.skillDescription, &res.effectName, &res.value, &res.duration, &res.tick)
		if err != nil {
			panic(err)
		}
		buildCharacter(*res, charMap)
	}

	return charMap
}

// BuildCharacter builds a character from our engine package
func buildCharacter(res dbFeed, charMap map[int]engine.Character) {
	_, ok := charMap[res.charID]
	key := res.charID
	// the follwing 2 ifs are here so we don't build the same character and skill more than once
	if !ok {
		charMap[key] = engine.Character{key, res.charName, 100, map[int]engine.Skill{}}
	}

	skills := charMap[key].Skills
	_, ok = skills[res.skillID]
	if !ok {
		skills[res.skillID] = engine.Skill{res.skillID, res.skillName, res.skillDescription, map[string][]engine.Effect{}}
	}

	effects := skills[res.skillID].Effects
	_, ok = effects[res.effectName]
	if !ok {
		effects[res.effectName] = make([]engine.Effect, 0)
	}

	if res.effectName == "damage" {
		effects[res.effectName] = append(effects[res.effectName], engine.Damage{res.value, res.tick, res.duration})
	}
}

func socket(conn *websocket.Conn) {
	defer conn.Close()
	type msg struct {
		Hello string `json:"hello"`
	}
	for {
		m := msg{}
		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading JSON.", err)
			break
		}

		fmt.Println(m)
		if m.Hello == "x" {
			fmt.Println("TERMINATING CONNECTION")
			break
		}
		if err = conn.WriteJSON(m); err != nil {
			fmt.Println("SOMETHING WENT WRONG")
		}
	}
}
