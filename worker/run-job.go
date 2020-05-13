package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func runJob(conn *websocket.Conn) {

	for i := 1; i < 5; i++ {
		time.Sleep(3 * time.Second)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("This is message # %d", i)))
	}

	conn.WriteMessage(websocket.CloseMessage, nil)
}
