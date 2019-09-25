package server

import (
	"fmt"
	"smash/gamedb"
)

func buildTeam(r *startGameReq) {

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
		fmt.Println(res)
	}
}
