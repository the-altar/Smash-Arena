package engine

// Skill is a struct for in-game abilities
type Skill struct {
	ID      int
	Name    string
	Desc    string
	Effects map[string][]Effect
}
