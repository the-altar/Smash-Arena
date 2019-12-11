package arena

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	json.Unmarshal(rawData, &p)

	if err != nil {
		return
	}

	p.Facepic = uuid.New().String()
	for _, skill := range p.Skills {
		skill.Skillpic = uuid.New().String()
	}

	fmt.Println(p.Facepic)
	fmt.Println("----------> CREATED")
	newPersona(p)
	g.Status(http.StatusOK)
}

// UpdatePersona updates a persona
func UpdatePersona(g *gin.Context) {
	fmt.Println("Uploading file")
	p := Persona{}
	j := g.Request.FormValue("json")
	fmt.Println(j)

	err := json.Unmarshal([]byte(j), &p)

	if err != nil {
		fmt.Println(err)
		g.Status(http.StatusInternalServerError)
		return
	}

	for index := range p.Skills {
		filename := fmt.Sprint("skillpic_", index)
		fmt.Println(filename + "++++++++")
		skillPicture, _ := g.FormFile(filename)

		if p.Skills[index].Skillpic == "" {
			fmt.Println(p.Skills[index].SkillName + " doesn't have a code!")
			p.Skills[index].Skillpic = uuid.New().String()
		}

		dst := "./public/img/character/skill/" + p.Skills[index].Skillpic + ".jpg"

		if skillPicture != nil {
			fmt.Println("File was found!")
			fmt.Println(skillPicture.Filename)
			g.SaveUploadedFile(skillPicture, dst)
		}
	}

	f, _ := g.FormFile("facepic")
	if f != nil {
		if p.Facepic == "null" || p.Facepic == "" {
			p.Facepic = uuid.New().String()
		}

		dst := "./public/img/character/profile/" + p.Facepic + ".jpg"
		g.SaveUploadedFile(f, dst)
	}
	updatePersona(p)
	g.Status(http.StatusOK)
}
