package arena

import (
	"encoding/json"
	"fmt"

	"github.com/the-altar/Smash-Arena/pkg/config"
)

// Persona struct
type Persona struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Profile  string `json:"profile"`
	Facepic  string `json:"facepic"`
	Skills   []struct {
		SkillName   string `json:"skillName"`
		Skillpic    string `json:"skillpic"`
		Description string `json:"description"`
		Cooldown    string `json:"cooldown"`
		Selection   int    `json:"selection"`
		Costs       []int  `json:"costs"`
		Effects     []struct {
			Type       int `json:"type"`
			Tick       int `json:"tick"`
			Duration   int `json:"duration"`
			Value      int `json:"value"`
			Condition  int `json:"condition"`
			Attr       int `json:"attr"`
			Trigger    int `json:"trigger"`
			AutoTarget int `json:"auto_target"`
		} `json:"effects"`
		Target int `json:"target"`
	} `json:"skills"`
}

// AllPersona from database
func allPersona() []Persona {
	p := make([]Persona, 0)

	sql := "SELECT * from persona"
	rows, err := config.DB.Query(sql)

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		p1 := Persona{}
		var s string
		if err = rows.Scan(&p1.ID, &p1.Nickname, &p1.Profile, &s, &p1.Facepic); err != nil {
			panic(err)
		}
		json.Unmarshal([]byte(s), &p1.Skills)
		p = append(p, p1)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return p
}

func newPersona(p Persona) {
	sql := "INSERT INTO public.persona (nickname, profile, gamedata, facepic) VALUES($1, $2, $3, $4);"
	rawData, err := json.Marshal(&p.Skills)

	if err != nil {
		fmt.Println(err)
		return
	}
	s := string(rawData)
	fmt.Println(s)
	if _, err = config.DB.Query(sql, p.Nickname, p.Profile, p.Facepic, s); err != nil {
		fmt.Println(err)
		return
	}
	return
}
