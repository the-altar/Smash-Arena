package server

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func matchmake() {
	for i := 1; i < rManager.poolSize(); i++ {
		_, g1 := rManager.poolPop()
		_, g2 := rManager.poolPop()

		g1.joinEnemy(g2.Game.GetTeam(), g2.Game.GetPlayer(), false)
		g2.joinEnemy(g1.Game.GetTeam(), g1.Game.GetPlayer(), true)

		rManager.GamePool.active[g1.Game.GetPlayer()] = g1
		rManager.GamePool.active[g2.Game.GetPlayer()] = g2
	}

	return
}

func listenSocket(g *gameHub, id string, chat chan int) {
	clientMsg := clientMessageGame{}
	timer := time.NewTicker(20 * time.Second)
	fanOut := make(chan bool)

	defer func() {
		g.ws.Close()
		timer.Stop()
	}()

	go func() {
		for {
			select {
			case <-timer.C:
				switchTurns(g)
			case <-fanOut:
				return
			}
		}
	}()

	for {
		if err := g.ws.ReadJSON(&clientMsg); err == nil {
			switch clientMsg.Code {
			case 0:
				timer.Stop()
				timer = time.NewTicker(20 * time.Second)
				switchTurns(g)
			case 1:
				fmt.Println("God this client is annoying...")
				g.send <- 1
			case 2:
				g.send <- 2
				fmt.Println("Sent channel!")
			case 3:
				fmt.Println(rManager.Rooms[id])
			}
		} else {
			chat <- 3
			fanOut <- true
			return
		}
	}
}

func serveSocket(g *gameHub, chat chan int) {
	messageGS := &clientMessageGame{}
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		g.ws.Close()
		ticker.Stop()
	}()

	for {
		select {
		case msg := <-g.send: // this ccl
			g.ws.WriteJSON(clientMessageGame{
				Client: "I received your message now stfu!",
				Code:   msg,
			})
		case status := <-g.ongoing: // this channel tells our server when two players have been matched'
			if status {
				fmt.Println("Working hard!")
				messageGS.writeGameState(g.Game)
				g.ws.WriteJSON(messageGS)
			}
		case <-chat: // chat tells me when the client has gone offline
			endGame(g)
			return
		case <-ticker.C:
			if err := g.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				chat <- 1
			}
		}
	}
}

func endGame(g *gameHub) {
	fmt.Println("Client left :(")
	rManager.removeFromPool(g.Game.GetPlayer())
	rManager.deleteRoom(g.Game.GetPlayer())

	delete(rManager.GamePool.active, g.Game.GetPlayer())

	fmt.Printf("Empty rooms left: %d\n", rManager.poolSize())
	fmt.Printf("Arenas remaining: %d\n", len(rManager.Rooms))
}

func switchTurns(g *gameHub) {
	fmt.Println("SWITCH! from: server_socket.go switchTurns() ")
	g.Game.YourTurn = !g.Game.YourTurn
	g.ongoing <- true
}
