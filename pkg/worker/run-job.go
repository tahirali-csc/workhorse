package worker

import (
	"bufio"
	"fmt"
	"workhorse/pkg/util"

	"github.com/gorilla/websocket"
)

func RunJob(conn *websocket.Conn) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error in reading message", err)
		return
	}

	//Deseriliaze the msg to an object
	job := util.ConvertToJobObject(msg)

	//Run the job in the container
	response := runDockerContainer(job)

	//Copy the conainer response to websocket stream
	rd := bufio.NewReader(response)

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		conn.WriteMessage(websocket.TextMessage, []byte(line))
	}

	conn.WriteMessage(websocket.CloseMessage, nil)
}
