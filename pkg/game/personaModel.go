package game

import (
	"encoding/json"

	"github.com/the-altar/Smash-Arena/pkg/config"
)

type Persona struct {
	Key      int      `json:key`
	Nickname string   `json:"nickname"`
	Profile  string   `json:profile`
	Health   int      `json:"health"`
	Skills   []Skills `json:"skills"`
}

type Skills struct {
	Skill      Skill     `json:"skill"`
	Params     Params    `json:"params"`
	JsonbAgg   []Effects `json:"jsonb_agg"`
	PersonaKey int       `json:"persona_key"`
}
type Skill struct {
	Cooldown    int    `json:"cooldown"`
	SkillKey    int    `json:"skill_key"`
	SkillName   string `json:"skill_name"`
	Description string `json:"description"`
	PersonaKey  int    `json:"persona_key"`
}
type Params struct {
	Hidden    int `json:"hidden"`
	Special   int `json:"special"`
	Physical  int `json:"physical"`
	Randomic  int `json:"randomic"`
	ParamKey  int `json:"param_key"`
	SkillKey  int `json:"skill_key"`
	Strategic int `json:"strategic"`
}
type Effects struct {
	Attr       string `json:"attr"`
	Tick       int    `json:"tick"`
	Value      int    `json:"value"`
	Target     int    `json:"target"`
	Duration   int    `json:"duration"`
	SkillKey   int    `json:"skill_key"`
	EffectKey  int    `json:"effect_key"`
	EffectType int    `json:"effect_type"`
}

//BuildTeam retrieves information about a character completely
func BuildTeam(ids []int) []Persona {
	sql := "select p.*, jsonb_agg(sel.*) as skills from persona as p join ( select s.persona_key, to_jsonb(s.*) as skill, to_jsonb(sc.*) as params, jsonb_agg(se.*) from skill as s join skill_costs as sc on s.skill_key = sc.skill_key join skill_effects as se on s.skill_key = se.skill_key where s.persona_key = $1 or s.persona_key = $2 or s.persona_key = $3 group by s.skill_key, sc.param_key, s.persona_key ) as sel on sel.persona_key = p.persona_key where p.persona_key = $1 or p.persona_key = $2 or p.persona_key = $3 group by p.persona_key;"
	personas := make([]Persona, 0)

	rows, err := config.DB.Query(sql, ids[0], ids[1], ids[2])
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		p := &Persona{Health: 100}
		var s json.RawMessage

		if err = rows.Scan(&p.Key, &p.Nickname, &p.Profile, &s); err != nil {
			panic(err)
		}
		json.Unmarshal(s, &p.Skills)
		personas = append(personas, *p)
	}

	return personas
}
