package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

var (
	db *sql.DB
)

type char struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type response struct {
	Roster []char `json:"roster"`
}

// InitServer boots our server
func InitServer(port string, database *sql.DB) {
	e := echo.New()
	db = database

	e.Static("/", "static")
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	e.GET("/characters", func(c echo.Context) error {
		var roster []char
		rows, err := db.Query("SELECT id, name FROM characters")
		if err != nil {
			return c.String(http.StatusOK, "Falhou")
		}
		defer rows.Close()

		for rows.Next() {
			var name, profile string
			if err := rows.Scan(&profile, &name); err != nil {
				return c.String(http.StatusOK, "Falhou feio")
			}
			roster = append(roster, char{name, profile})
		}

		U := &response{roster}
		fmt.Println(U)
		return c.JSON(http.StatusOK, U)
	})

	/*var engine engine.GameRoom
	team := [3]string{"N", "R", "T"}
	engine.Begin("Joao", team)*/
	e.Logger.Fatal(e.Start(":" + port))
}
