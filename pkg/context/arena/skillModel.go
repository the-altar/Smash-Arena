package arena

import (
	"encoding/json"
	"fmt"

	"github.com/the-altar/Smash-Arena/pkg/config"
)

type main struct {
	SkillKey    int    `json:"skill_key"`
	SkillName   string `json:"skill_name"`
	Description string `json:"description"`
	PersonaKey  int    `json:"persona_key"`
	Cooldown    int    `json:"cooldown"`
}

type params struct {
	ParamKey  int `json:"param_key"`
	Physical  int `json:"physical"`
	Special   int `json:"special"`
	Hidden    int `json:"hidden"`
	Strategic int `json:"strategic"`
	Randomic  int `json:"randomic"`
	SkillKey  int `json:"skill_key"`
}

type skill struct {
	Main   main
	Params params
}

func oneSkillSet(charID int) []skill {
	skills := make([]skill, 0)
	sql := "select to_jsonb(s.*) as skill, to_jsonb(sc.*) as params from skill as s join skill_costs as sc on s.skill_key = sc.skill_key where s.persona_key = $1 group by s.skill_key, sc.param_key;"

	rows, err := config.DB.Query(sql, charID)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		s := skill{}
		var r1, r2 json.RawMessage

		if err = rows.Scan(&r1, &r2); err != nil {
			panic(err)
		}
		json.Unmarshal(r1, &s.Main)
		json.Unmarshal(r2, &s.Params)

		skills = append(skills, s)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	fmt.Println(skills)
	return skills
}
