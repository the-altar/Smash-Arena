package main

import (
	"os"
	"smash/database"
	"smash/server"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	connString := os.Getenv("DATABASE_URL")
	database := database.RunDB(connString)

	if port == "" {
		port = "3000"
	}
	server.InitServer(port, database)
}
