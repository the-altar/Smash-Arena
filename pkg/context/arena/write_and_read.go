package arena

import (
	"fmt"
	"time"

	"github.com/the-altar/Smash-Arena/pkg/game"
)

func closeDoor(gid string) {
	r := Conn.rooms[gid]
	if len(r.player) == 0 {
		delete(Conn.rooms, gid)
		fmt.Printf("ROOM %s is closed!\n", gid)
		r.isdestroyed <- true
	}
}

func serveRoom(r *rooms) {
	for {
		select {
		case <-r.turn:
			r.switchPlayerTurn()

		case index := <-r.playerleft:
			r.player[index].isdestroyed <- r.gid
			go Conn.clearConn(r.player[index].pid, r.gid, index)

		case <-time.After(60 * time.Second):
			r.switchPlayerTurn()

		case <-r.isdestroyed:
			fmt.Println("\n**** Showtime is over! ****")
			return
		}
	}
}

func readPump(c *connection) {

	r := &ClientMessage{}
	for {
		err := c.client.ReadJSON(r)
		fmt.Println(r)
		if err != nil {
			c.client.Close()

			if c.gamePos < 0 {
				c.isdestroyed <- ""
				Conn.clearConn(c.pid, "", 0)
				return
			}

			fmt.Println("Player left, waiting for reconnect")
			select {
			case <-time.After(60 * time.Second):
				Conn.rooms[c.gid].playerleft <- c.gamePos
				return
			case <-c.reconnect:
				fmt.Println("player reconnected")
				c.isready <- 3
			}
		}
		if r.Code == 1 {
			fmt.Println("Client connected and solicited that I build his team :)")
			if c.Arena != nil {
				fmt.Println("But they're already ingame :(")
			} else {
				c.Arena = &game.Arena{}
				c.Arena.Allies = game.BuildTeam(r.TeamID)
			}

		} else if r.Code == 2 {
			if Conn.rooms[c.gid].isItYourTurn(c.pid) {
				Conn.rooms[c.gid].turn <- r.Code
			}
		}
	}
}

func writePump(c *connection) {
	for {
		select {
		case code := <-c.isready:
			if code == 1 {
				msg := &ClientMessage{
					Code:     0,
					Gid:      c.gid,
					GameData: *c.Arena,
				}
				fmt.Printf("Hey %s, it's your turn \n", c.pid)
				c.client.WriteJSON(msg)
			} else if code == 3 {
				c.client.WriteJSON(&ClientMessage{
					Code:     3,
					Gid:      c.gid,
					GameData: *c.Arena,
				})
				fmt.Println("heh, so you're back. Here's your game data")
			}
		case <-c.isdestroyed:
			return

		case <-time.After(50 * time.Second):
			c.client.WriteJSON("ping")
		}
	}
}
