package engine

// Character represent a char which the player controls
type character struct {
	health int
	name   string
	alive  bool
}

// Build initializes a character for the game
func (c *character) build(name string) {
	c.health = 100
	c.name = name
	c.alive = true
}

func (c *character) takeDamage(damage int) {
	if c.health > 0 {
		c.health = c.health - damage
	} else {
		return
	}

	if c.health <= 0 {
		c.alive = false
	}
}
