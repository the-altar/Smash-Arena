package engine

// Skill is a struct for in-game abilities
type Skill struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	effects map[string][]Effect
}

// BuildSkill builds, that's right you guessed it, a new skill!
func BuildSkill(id int, name string, desc string) Skill {
	s := Skill{id, name, desc, make(map[string][]Effect)}

	return s
}

// GetEffects returns an slice of an effect type through a parameter
func (s Skill) GetEffects() map[string][]Effect {
	return s.effects
}
