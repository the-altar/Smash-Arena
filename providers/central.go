package providers

import (
	"fmt"
	"time"
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

func readPump(c *connection, counter int, lastResponse time.Time) {

	r := &ClientMessage{}
	for {
		err := c.client.ReadJSON(r)
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
			}
		}
		if r.Code == 2 {
			if Conn.rooms[c.gid].isItYourTurn(c.pid) {
				Conn.rooms[c.gid].turn <- r.Code
			}
		}
	}
}

func writePump(c *connection, counter int, lastResponse time.Time) {
	fmt.Println("Write pumping...")
	for {
		select {
		case <-c.isready:
			msg := &ClientMessage{
				Code: 0,
				Gid:  c.gid,
				Data: Data{Test: 1},
			}
			fmt.Printf("Hey %s, it's your turn dude", c.pid)
			c.client.WriteJSON(msg)

		case <-c.isdestroyed:
			return

		case <-time.After(50 * time.Second):
			c.client.WriteJSON("ping")
		}
	}
}
