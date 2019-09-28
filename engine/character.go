package engine

// Character defines a playable character
type Character struct {
	ID     int
	Name   string
	Health int
	Skills map[int]Skill
}

func (c Character) getSkills(string) map[int]Skill {
	return c.Skills
}
