package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func socket(ws *websocket.Conn) {
	clientMsg := clientMessageGame{}
	defer ws.Close()

	for {
		if err := ws.ReadJSON(&clientMsg); err == nil {
			switch clientMsg.Code {
			case 1:
				break
			case 2:
				fmt.Println(clientMsg.Data)
			case 3:
				fmt.Println(arenas[clientMsg.Client])
			}
		} else {
			fmt.Println(err)
		}
	}
}
