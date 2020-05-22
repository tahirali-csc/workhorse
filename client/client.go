package main

import (
	"fmt"
	"log"
	"net/url"
	"workhorse/util"

	"github.com/gorilla/websocket"
)

func main() {

	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/runWorkflow"}
	// u := url.URL{Scheme: "ws", Host: "192.168.56.102:8081", Path: "/runWorkflow"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	workflowBytes := util.ConvertToByteArray(readWorkflow())
	conn.WriteMessage(websocket.BinaryMessage, workflowBytes)

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
