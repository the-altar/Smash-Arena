package gamedb

import "fmt"

// QueryCharData buils a query string for a specific skills
func QueryCharData(TeamID []string) string {

	arrSize := len(TeamID)
	query := `select c.char_id, c.char_name, s.skill_id, s.skill_name, s.skill_description, e.effect_name, ep.value, ep.duration, ep.tick from "characters" as c join skills as s on s.char_id = c.char_id join effects as e on e.skill_id = s.skill_id join effect_params as ep on ep.effect_id = e.effect_id where`

	for i := 0; i < arrSize; i++ {
		qString := " c.char_id = " + TeamID[i]

		if i < (arrSize - 1) {
			qString = qString + " or"
		}

		query = query + qString
	}
	fmt.Println(query)
	return query
}
