package server

import (
	"net/http"
	"smash/engine"
	"smash/gamedb"

	"github.com/labstack/echo"
)

func startGameHandler(c echo.Context) error {
	var game map[int]engine.Character

	r := &startGameReq{}

	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, 0)
	}
	game = buildTeam(r) // from ./server_helpers.go
	return c.JSON(http.StatusOK, game)
}

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

func buildCharacter(res dbFeed, charMap map[int]engine.Character) {
	_, ok := charMap[res.charID]
	key := res.charID
	if !ok {
		charMap[key] = engine.Character{res.charID, res.charName, 100, map[int]engine.Skill{}}
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
