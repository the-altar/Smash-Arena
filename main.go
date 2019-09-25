package main

import (
	"os"
	"smash/gamedb"
	"smash/server"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	connString := os.Getenv("DATABASE_URL")
	database := gamedb.RunDB(connString)

	if port == "" {
		port = "3000"
	}
	server.InitServer(port, database)
}
