package engine

// Game defines a struct for a game room
type Game struct {
	party party
	full  bool
}

// Begin starts a game
func (g *Game) Begin(player string, members [3]string) {
	g.party.form(members, player)
}
