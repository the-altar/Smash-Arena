package database

import (
	"database/sql"
	"log"
)

// RunDB kicks starts our database
func RunDB(conn string) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
