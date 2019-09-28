package server

// startGameReq is the initial request we get from the client to start a game
type startGameReq struct {
	UserID string   `json:"UserID"`
	TeamID []string `json:"TeamID"`
}

/* dbFeed contains all the info we get from the database. I'm using it as an struct to avoid passing
a billion parameters to the constructor function*/
type dbFeed struct {
	charID           int
	charName         string
	skillID          int
	skillName        string
	skillDescription string
	effectName       string
	value            int
	duration         int
	tick             bool
}

// charClient is what we send back to the client when they request information about a character
type charClient struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Profile string `json:"Profile"`
}
