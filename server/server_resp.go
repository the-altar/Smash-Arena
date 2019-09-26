package server

type startGameReq struct {
	UserID string   `json:"UserID"`
	TeamID []string `json:"TeamID"`
}
