package server

import (
	"net/http"
	"smash/gamedb"

	"github.com/labstack/echo"
)

func getCharactersHandler(c echo.Context) error {
	type char struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	type response struct {
		Roster []char `json:"roster"`
	}

	roster := &response{make([]char, 0)}
	query := gamedb.QueryAllCharShallow()
	rows, err := db.Query(query)

	if err != nil {
		return c.String(http.StatusOK, "Falhou")
	}

	defer rows.Close()

	for rows.Next() {
		char := char{}
		if err := rows.Scan(&char.ID, &char.Name); err != nil {
			return c.String(http.StatusOK, "Falhou feio")
		}
		roster.Roster = append(roster.Roster, char)
	}

	return c.JSON(http.StatusOK, roster)
}
