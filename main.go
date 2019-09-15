package main

import (
	"smash/engine"
	"smash/server"
)

func main() {
	var game engine.Game
	m := [3]string{"Sasuke", "Sakura", "Naruto"}
	game.Begin("Joao", m)

	server.InitServer()
}
