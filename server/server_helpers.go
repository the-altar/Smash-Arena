package server

import (
	"fmt"
	"smash/engine"
	"smash/gamedb"
)

type results struct {
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

func buildTeam(r *startGameReq) {
	charMap := make(map[int]engine.Character)
	teamQuery := gamedb.QueryCharData(r.TeamID)
	rows, err := db.Query(teamQuery)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		res := &results{}
		err = rows.Scan(&res.charID, &res.charName, &res.skillID, &res.skillName, &res.skillDescription, &res.effectName, &res.value, &res.duration, &res.tick)
		if err != nil {
			panic(err)
		}
		buildCharacter(*res, charMap)
	}
	fmt.Println(charMap)
}

func buildCharacter(res results, charMap map[int]engine.Character) {
	_, ok := charMap[res.charID]
	key := res.charID
	if !ok {
		charMap[key] = engine.Character{res.charID, res.charName, 100, map[int]engine.Skill{}}
	}

	_, ok = charMap[key].Skills[res.skillID]
	if !ok {
		charMap[key].Skills[res.skillID] = engine.Skill{res.skillID, res.skillName, res.skillDescription, map[string][]engine.Effect{}}
	}

	_, ok = charMap[key].Skills[res.skillID].Effects[res.effectName]
	if !ok {
		charMap[key].Skills[res.skillID].Effects[res.effectName] = make([]engine.Effect, 0)
	}

	charMap[key].Skills[res.skillID].Effects[res.effectName] = append(charMap[key].Skills[res.skillID].Effects[res.effectName], engine.Damage{res.value, res.tick, res.duration})
}
