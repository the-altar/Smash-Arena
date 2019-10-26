package config

import (
	"database/sql"
	"fmt"
	"os"

	// local driver
	_ "github.com/lib/pq"
)

// DB is the database variable
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
