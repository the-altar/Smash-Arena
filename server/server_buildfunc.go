package server

import (
	"smash/engine"
	"smash/gamedb"
)

// BuildTeam builds the player's team for the game
func buildTeam(r *startGameReq) map[string]engine.Character {
	charMap := make(map[string]engine.Character)
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

func buildGameRoom(r *startGameReq) engine.GameRoom {
	team := buildTeam(r)
	gRoom := engine.GameRoom{}
	gRoom.AddTeam(team)
	gRoom.SetTimer(60)
	gRoom.SetPlayer(r.UserID)
	return gRoom
}

// BuildCharacter builds a character from our engine package
func buildCharacter(res dbFeed, charMap map[string]engine.Character) {
	_, ok := charMap[res.charName]
	key := res.charName
	// the follwing 2 ifs are here so we don't build the same character and skill more than once
	if !ok {
		charMap[res.charName] = engine.Character{res.charID, res.charName, 100, map[string]engine.Skill{}}
	}

	skills := charMap[key].Skills
	_, ok = skills[res.skillName]
	if !ok {
		skills[res.skillName] = engine.BuildSkill(res.skillID, res.skillName, res.skillDescription)
	}

	effects := skills[res.skillName].GetEffects()
	_, ok = effects[res.effectName]
	if !ok {
		effects[res.effectName] = make([]engine.Effect, 0)
	}

	if res.effectName == "damage" {
		effects[res.effectName] = append(effects[res.effectName], engine.Damage{res.value, res.tick, res.duration})
	}
}
