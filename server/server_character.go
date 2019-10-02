package server

import (
	"net/http"
	"smash/gamedb"

	"github.com/labstack/echo"
)

func getCharactersHandler(c echo.Context) error {

	type response struct {
		Roster []charClient `json:"roster"`
	}

	roster := &response{make([]charClient, 0)}
	query := gamedb.QueryAllCharShallow()
	rows, err := db.Query(query)

	if err != nil {
		return c.String(http.StatusOK, "Falhou")
	}

	defer rows.Close()

	for rows.Next() {
		char := charClient{}
		if err := rows.Scan(&char.ID, &char.Name, &char.Profile); err != nil {
			return c.String(http.StatusOK, "Falhou feio")
		}
		roster.Roster = append(roster.Roster, char)
	}

	return c.JSON(http.StatusOK, roster)
}
