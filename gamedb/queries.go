package gamedb

// QuerySkillData buils a query string for a specific skills
func QuerySkillData(skillID []string) string {
	query := "select s.skill_name, 	e.effect_name, 	ep.value from skills as s join effects as e	on e.skill_id = s.skill_id join effect_params as ep	on ep.effect_id = e.effect_id where "

	for i := 0; i < len(skillID); i++ {
		query = query + "s.skill_id = " + skillID[i]
	}

	query = query + ";"
	return query
}
