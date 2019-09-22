package server

type char struct {
	Name    string  `json:"name"`
	Profile string  `json:"profile"`
	Skills  []skill `json:"skills"`
}

type skill struct {
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
}

type response struct {
	Roster []char `json:"roster"`
}
