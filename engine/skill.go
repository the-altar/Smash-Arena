package engine

// Skill is a struct for in-game abilities
type Skill struct {
	ID      int
	Name    string
	Desc    string
	Effects map[string][]Effect
}

// GetEffect returns an slice of an effect type through a parameter
func (s Skill) GetEffect(effectType string) []Effect {
	return s.Effects[effectType]
}
