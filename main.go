package main

import (
	"os"
	"smash/server"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}
	server.InitServer(port)
}
