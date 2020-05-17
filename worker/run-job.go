package main

import (
	"bufio"

	"github.com/gorilla/websocket"
)

func runJob(conn *websocket.Conn) {

	// for i := 1; i < 5; i++ {
	// 	time.Sleep(3 * time.Second)
	// 	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("This is message # %d", i)))
	// }

	response := runDockerContainer()
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
