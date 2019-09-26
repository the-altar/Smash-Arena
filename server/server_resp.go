package server

type startGameReq struct {
	UserID string   `json:"UserID"`
	TeamID []string `json:"TeamID"`
}

type char struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Profile string `json:"Profile"`
}
