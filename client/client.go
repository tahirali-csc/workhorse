package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	//u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/runWorkflow"}
	u := url.URL{Scheme: "ws", Host: "192.168.56.102:8081", Path: "/runWorkflow"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	conn.WriteMessage(websocket.BinaryMessage, []byte("Hello Client"))

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.CloseMessage {
			break
		} else {
			fmt.Print(string(msg))
		}
	}

	defer conn.Close()
}
