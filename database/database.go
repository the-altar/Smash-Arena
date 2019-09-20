package database

import (
	"database/sql"
	"fmt"
	"log"
)

//RunDB kicks starts our database
func RunDB(conn string) *sql.DB {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database!")
	return db
}

//CloseDB terminates the connection
func CloseDB(db *sql.DB) {
	db.Close()
}
