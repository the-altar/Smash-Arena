package server

import (
	"fmt"
	"time"
)

func matchmake() {
	for i := 1; i < rManager.poolSize(); i++ {
		_, g1 := rManager.poolPop()
		_, g2 := rManager.poolPop()

		g1.joinEnemy(g2.Game.GetTeam(), g2.Game.GetPlayer())
		g2.joinEnemy(g1.Game.GetTeam(), g1.Game.GetPlayer())
	}
	return
}

func listenSocket(g *gameHub, id string, chat chan int, t *time.Time) {
	clientMsg := clientMessageGame{}

	defer g.ws.Close()
	for {
		if err := g.ws.ReadJSON(&clientMsg); err == nil {
			switch clientMsg.Code {
			case 0:
				*t = time.Now()
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
			return
		}
	}
}

func serveSocket(g *gameHub, chat chan int, t *time.Time) {
	defer g.ws.Close()
	ticker := time.NewTicker(50 * time.Second)
	messageGS := &clientMessageGame{}

	for {
		select {
		case msg := <-g.send: // this ccl
			g.ws.WriteJSON(clientMessageGame{
				Client: "I received your message now stfu!",
				Code:   msg,
			})
		case <-g.ongoing: // this channel tells our server when two players have been matched'
			messageGS.writeGameState(g.Game)
			g.ws.WriteJSON(messageGS)

		case <-chat: // chat tells me when the client has gone offline
			fmt.Println("Client left :(")
			rManager.removeFromPool(g.Game.GetPlayer())
			rManager.deleteRoom(g.Game.GetPlayer())

			fmt.Printf("Empty rooms left: %d\n", rManager.poolSize())
			fmt.Printf("Arenas remaining: %d\n", len(rManager.Rooms))
			return
		case <-ticker.C:
			go func() {
				g.ws.WriteJSON("pong")
				if time.Now().Sub(*t) > 100 {
					chat <- 1
					ticker.Stop()
				}
			}()
		}
	}
}
