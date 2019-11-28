package arena

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GameSocket will handle our websocket connection
func GameSocket(g *gin.Context) {
	v := make(chan bool)
	go Conn.Begin(g, v)

	Conn.PumpOut(g.Param("id"), <-v)
}

// Arena serves our SPA
func Arena(g *gin.Context) {
	g.File("./public/arena/index.html")
	return
}

// AllPersona Fetches all personas from DB
func AllPersona(g *gin.Context) {
	personas := allPersona()
	g.JSON(http.StatusOK, personas)
	return
}

// OneSkillSet from a character
func OneSkillSet(g *gin.Context) {
	_, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, 0)
	}

	//s := oneSkillSet(id)
	g.JSON(http.StatusOK, nil)
	return
}

// CreatePersona adds a new character to the database
func CreatePersona(g *gin.Context) {
	p := Persona{}
	rawData, err := g.GetRawData()
	if err != nil {
		return
	}

	json.Unmarshal(rawData, &p)

	newPersona(p)
}
