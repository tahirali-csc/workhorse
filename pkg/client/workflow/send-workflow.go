package workflow

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func SendWorkFlow(masterNodeAddress string, workflowTransferObjBytes []byte) {
	u := url.URL{Scheme: "ws", Host: masterNodeAddress, Path: "/runWorkflow"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	log.Println("Sending the workflow to server :::" + masterNodeAddress)
	conn.WriteMessage(websocket.BinaryMessage, workflowTransferObjBytes)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.CloseMessage {
			break
		} else {
			//Write to console
			fmt.Print(string(msg))
		}
	}

	defer conn.Close()
}
