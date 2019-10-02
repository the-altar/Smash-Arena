package engine

// GameRoom defines a game room for each player
type GameRoom struct {
	opponent string
	player   string
	yourturn bool
	Full     chan bool
	timer    int
	team     map[int]Character
	enemies  map[int]Character
}

// AddTeam adds characters to our "team" field
func (g *GameRoom) AddTeam(chars map[int]Character) {
	g.team = chars
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
	g.Full = make(chan bool)
	g.player = player
}

// SetOpponent sets player 2
func (g *GameRoom) SetOpponent(player string) {
	g.opponent = player
}

// GetPlayer returns the player's id
func (g *GameRoom) GetPlayer() string {
	return g.player
}

// GetTeam returns the player's team
func (g *GameRoom) GetTeam() map[int]Character {
	return g.team
}
