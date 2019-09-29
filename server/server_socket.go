package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type clientMessage struct {
	client string
	code   int
	data   interface{}
}

func socket(ws *websocket.Conn) {
	clientMsg := &clientMessage{}
	defer ws.Close()

	for {
		ws.ReadJSON(clientMsg)
		switch clientMsg.code {
		case 1:
			break
		case 2:
			println(clientMsg.data)
		case 3:
			println(arenas[clientMsg.client])
		default:
			fmt.Println("NOTHING TO SEE HERE")
		}
	}
}
