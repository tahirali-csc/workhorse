package workflow

import (
	"log"
	"net/url"
	"workhorse/pkg/api"
	"workhorse/pkg/util"

	"github.com/gorilla/websocket"
)

type Executor struct {
}

func (exec *Executor) Execute(masterNodeAddress string, workflow *api.Workflow) {

	//Convert workflow to byte array
	workflowByteObj, err := util.ConvertToByteArray(*workflow)
	if err != nil {
		panic(err)
	}

	//TODO: Will review this approach using using Websocket
	u := url.URL{Scheme: "ws", Host: masterNodeAddress, Path: "/runWorkflow"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("Unable to close connection. Error: %v", err)
		} else {
			log.Print("Connection closed")
		}
	}()

	log.Println("Sending the workflow to server :::" + masterNodeAddress)
	err = conn.WriteMessage(websocket.BinaryMessage, workflowByteObj)
	if err != nil {
		panic(err)
	}

	for {
		msgType, _, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.CloseMessage {
			break
		}
	}
}
