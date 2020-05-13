package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleWorkFlow(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		panic(err)
	}

	runWorkFlow(conn)

	defer func() {
		fmt.Println("Closing socket connection")
		conn.Close()
	}()
}

func main() {
	http.HandleFunc("/runWorkflow", handleWorkFlow)

	const addr = "localhost:8081"
	fmt.Println("Starting master node at:::" + addr)
	http.ListenAndServe(addr, nil)
}
