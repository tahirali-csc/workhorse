package main

import (
	"fmt"
	"net/http"
	"workhorse/pkg/util"
	"workhorse/pkg/worker"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleJob(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		panic(err)
	}

	worker.RunJob(conn)

	defer func() {
		fmt.Println("Closing socket connection")
		conn.Close()
	}()
}

func main() {
	http.HandleFunc("/runJob", handleJob)
	ipAddress := util.GetHostIPAddress()

	http.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "Hello from worker::: %s!", ipAddress)
	})

	addr := ":8080"
	fmt.Println("Starting worker node at:::" + addr)
	http.ListenAndServe(addr, nil)
}
