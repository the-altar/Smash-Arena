package engine

// GameRoom defines a struct for a game room
type GameRoom struct {
	party party
}

// Begin starts a game
func (g *GameRoom) Begin(player string, members [3]string) {
	g.party.form(members, player)
	var d damage
	d.construct(30, nil)
	var skill skill
	e := []effect{d}
	skill.build(1, "sdas", e)
}
