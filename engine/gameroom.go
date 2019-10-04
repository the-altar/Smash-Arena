package engine

// GameRoom defines a game room for each Player
type GameRoom struct {
	Opponent string               `json:"opponent"`
	Player   string               `json:"player"`
	YourTurn bool                 `json:"yourturn"`
	Timer    int                  `json:"timer"`
	Team     map[string]Character `json:"team"`
	Enemies  map[string]Character `json:"enemies"`
}

// AddTeam adds characters to our "Team" field
func (g *GameRoom) AddTeam(chars map[string]Character) {
	g.Team = chars
}

// AddEnemies adds characters to our "Enemies" field
func (g *GameRoom) AddEnemies(chars map[string]Character) {
	g.Enemies = chars
}

// SetTimer sets a Timer for the Player's turn
func (g *GameRoom) SetTimer(time int) {
	g.Timer = time
}

// SetPlayer sets Player 1
func (g *GameRoom) SetPlayer(Player string) {
	g.Player = Player
}

// SetOpponent sets Player 2
func (g *GameRoom) SetOpponent(Player string) {
	g.Opponent = Player
}

// GetPlayer returns the Player's id
func (g *GameRoom) GetPlayer() string {
	return g.Player
}

// GetTeam returns the Player's Team
func (g *GameRoom) GetTeam() map[string]Character {
	return g.Team
}
