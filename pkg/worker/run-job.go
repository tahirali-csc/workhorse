package worker

import (
	"bufio"
	"log"
	"workhorse/pkg/util"

	"github.com/gorilla/websocket"
)

func RunJob(conn *websocket.Conn) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Error in reading message", err)
		return
	}

	//Deseriliaze the msg to an object
	job := util.ConvertToJobObject(msg)

	//Run the job in the container
	response := runDockerContainer(job)

	//Copy the conainer response to websocket stream
	rd := bufio.NewReader(response)

	for {
		//TODO: Will review the logic again. Currently stream container logs line by line
		line, err := rd.ReadBytes('\n')
		if err != nil {
			break
		}

		conn.WriteMessage(websocket.BinaryMessage, line)
	}

	defer func() {
		log.Print("Sending close message to web socket")
		conn.WriteMessage(websocket.CloseMessage, nil)
	}()
}
