package server

import (
	"fmt"
)

func matchmake() {
	mutex.Lock()
	for i := 1; i < len(gamePool); i++ {
		g1 := gamePool[i-1]
		g2 := gamePool[i]
		gamePool = gamePool[i:]

		g1.game.AddEnemies(g2.game.GetTeam())
		g1.game.SetOpponent(g2.game.GetPlayer())

		g2.game.AddEnemies(g1.game.GetTeam())
		g2.game.SetOpponent(g1.game.GetPlayer())

		g2.game.Full <- true
		g1.game.Full <- true
	}
	mutex.Unlock()
}

func listenSocket(g gameHub, id string, chat chan int) {
	clientMsg := clientMessageGame{}

	defer g.ws.Close()
	for {
		if err := g.ws.ReadJSON(&clientMsg); err == nil {
			switch clientMsg.Code {
			case 1:
				fmt.Println("God this client is annoying...")
				g.send <- 1
			case 2:
				g.send <- 2
				fmt.Println("Sent channel!")
			case 3:
				fmt.Println(arenas[id])
			}
		} else {
			delete(arenas, id)
			chat <- 1
			break
		}
	}
}

func serveSocket(g gameHub, chat chan int) {
	defer g.ws.Close()
	for {
		select {
		case msg := <-g.send: // this channel tells me when a message is received from the client
			fmt.Println(msg)
			g.ws.WriteJSON(clientMessageGame{
				Client: "I received your message now stfu!",
				Code:   msg,
			})
		case isFull := <-g.game.Full: // the Full channel tells our code when the client's gameroom has found a match'
			if isFull {
				g.ws.WriteJSON(clientMessageGame{
					Client: "You've joined a room!",
					Code:   1,
				})
			} else {
				g.ws.WriteJSON(clientMessageGame{
					Client: "Nope, still waiting...",
					Code:   0,
				})
			}
		case <-chat: // chat tells me when the client has gone offline
			break
		}
	}
	fmt.Println("Listener is closed")
}
