package engine

// Character defines a playable character
type Character struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Health int              `json:"health"`
	Skills map[string]Skill `json:"skills"`
}

func (c Character) getSkills(string) map[string]Skill {
	return c.Skills
}
