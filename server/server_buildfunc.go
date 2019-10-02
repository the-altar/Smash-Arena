package server

import (
	"smash/engine"
	"smash/gamedb"
)

func buildGameRoom(charMap map[int]engine.Character, user string) *engine.GameRoom {
	gRoom := &engine.GameRoom{}
	gRoom.AddTeam(charMap)
	gRoom.SetTimer(60)
	gRoom.SetPlayer(user)
	return gRoom
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
