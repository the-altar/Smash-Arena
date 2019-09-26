package engine

// Character defines a playable character
type Character struct {
	ID     int
	Name   string
	Health int
	Skills map[int]Skill
}
