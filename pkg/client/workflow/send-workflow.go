package workflow

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func SendWorkFlow(masterNodeIP string, workflowTransferObjBytes []byte) {
	u := url.URL{Scheme: "ws", Host: masterNodeIP + ":8081", Path: "/runWorkflow"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending the workflow to server :::" + masterNodeIP)
	conn.WriteMessage(websocket.BinaryMessage, workflowTransferObjBytes)

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
