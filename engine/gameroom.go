package engine

// GameRoom defines a game room for each player
type GameRoom struct {
	opponent string
	player   string
	yourturn bool
	timer    int
	servant  map[int]Character
	enemies  map[int]Character
}

// AddServant adds characters to our "servant" field
func (g *GameRoom) AddServant(chars map[int]Character) {
	g.servant = chars
}

// AddEnemies adds characters to our "enemies" field
func (g *GameRoom) AddEnemies(chars map[int]Character) {
	g.enemies = chars
}

// SetTimer sets a timer for the player's turn
func (g *GameRoom) SetTimer(time int) {
	g.timer = time
}

// SetPlayer sets player 1
func (g *GameRoom) SetPlayer(player string) {
	g.player = player
}

// SetOpponent sets player 2
func (g *GameRoom) SetOpponent(player string) {
	g.player = player
}
