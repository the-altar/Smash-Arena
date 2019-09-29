package engine

import "fmt"

// GameRoom defines a game room for each player
type GameRoom struct {
	servant map[int]Character
	enemies map[int]Character
}

// AddServant adds characters to our "servant" field
func (g *GameRoom) AddServant(chars map[int]Character) {
	g.servant = chars
	fmt.Println(g.servant)
}

// AddEnemies adds characters to our "enemies" field
func (g *GameRoom) AddEnemies(chars map[int]Character) {
	g.enemies = chars
}
