package server

import (
	"fmt"
	"time"
)

func joinRoom(g gameHub, id string) {
	freeSize := len(freeArenas)
	g.game = arenas[id]

	for {
		if freeSize > 0 {
			opponent := freeArenas[0]
			freeArenas = freeArenas[:freeSize-1]

			opponent.AddEnemies(g.game.GetTeam())
			opponent.SetOpponent(g.game.GetPlayer())

			g.game.AddEnemies(opponent.GetTeam())
			g.game.SetOpponent(opponent.GetPlayer())

			g.game.Full <- true
			opponent.Full <- true
			break

		} else {
			freeArenas = append(freeArenas, g.game)
			g.game.Full <- false
			time.Sleep(5 * time.Second)
		}
	}

}

func listenSocket(g gameHub, id string) {
	clientMsg := clientMessageGame{}
	defer g.ws.Close()

	go joinRoom(g, id)

	for {
		if err := g.ws.ReadJSON(&clientMsg); err == nil {
			switch clientMsg.Code {
			case 1:

				fmt.Println("Game sucessfully built")
			case 2:
				g.send <- 1
				fmt.Println("Sent channel!")
			case 3:
				fmt.Println(arenas[id])
			}
		} else {
			delete(arenas, id)
			break
		}
	}
}

func serveSocket(g gameHub) {
	defer g.ws.Close()
	for {
		select {
		case msg := <-g.send:
			fmt.Println(msg)
			g.ws.WriteJSON(clientMessageGame{
				Client: "Server",
				Code:   2,
			})
		case isFull := <-g.game.Full:
			if isFull {
				g.ws.WriteJSON(clientMessageGame{
					Client: "ROOM IS FULL",
					Code:   1,
				})
			} else {
				g.ws.WriteJSON(clientMessageGame{
					Client: "Searching for an opponent! ",
					Code:   0,
				})
			}
		}
	}
}
