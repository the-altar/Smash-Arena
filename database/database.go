package database

import (
	"database/sql"
	"fmt"
	"log"
)

// RunDB kicks starts our database
func RunDB(conn string) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

}
