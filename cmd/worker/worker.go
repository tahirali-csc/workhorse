package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"workhorse/pkg/api"
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

	defer conn.Close()
}

func main() {
	serverAddr := flag.String("serverAddress", "localhost", "")
	serverPort := flag.Uint("serverPort", 8080, "")
	flag.Parse()

	serverInfo := &api.ServerConfig{
		Address: *serverAddr,
		Port:    *serverPort,
	}
	go worker.KeepSendingStats(serverInfo)

	http.HandleFunc("/runJob", handleJob)

	ipAddress := util.GetHostIPAddress()
	http.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {
		_ = util.GetSenderIP(request)
		fmt.Fprintf(w, "Hello from worker::: %s!", ipAddress)
	})

	addr := ":8080"
	log.Println("Starting worker node at:::" + addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Print(err)
	}
}
