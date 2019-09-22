package server

import (
	"net/http"

	"github.com/labstack/echo"
)

func getCharactersHandler(c echo.Context) error {
	var roster []char
	query := `select 	c.id, c.name, s.id as skill_id, s."name" as skill_id, s.description from "characters" as c join skills as s on c.id = s.char_id group by c.id, s.id;`

	rows, err := db.Query(query)
	if err != nil {
		return c.String(http.StatusOK, "Falhou")
	}

	defer rows.Close()

	for rows.Next() {
		var name, profile, skillID, skillName, description string
		var skills []skill
		if err := rows.Scan(&profile, &name, &skillID, &skillName, &description); err != nil {
			return c.String(http.StatusOK, "Falhou feio")
		}
		skills = append(skills, skill{skillName, skillID, description})
		roster = append(roster, char{name, profile, skills})
	}

	U := &response{roster}
	return c.JSON(http.StatusOK, U)
}
