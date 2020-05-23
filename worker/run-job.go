package main

import (
	"bufio"
	"fmt"
	"workhorse/util"

	"github.com/gorilla/websocket"
)

func runJob(conn *websocket.Conn) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error in reading message", err)
		return
	}

	//Deseriliaze the msg to an object
	job := util.ConvertToJobObject(msg)

	response := runDockerContainer(job)
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
