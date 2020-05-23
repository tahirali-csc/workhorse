package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"workhorse/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var WorkScheduler Scheduler = nil

func handleWorkFlow(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		panic(err)
	}

	runWorkFlowSync(conn, WorkScheduler)

	defer func() {
		fmt.Println("Closing socket connection")
		conn.Close()
	}()
}

func main() {
	//Read command line arguments
	workerNodeAddrParams := flag.String("worker-node-address", "", "Comma separted ip address of worker nodes")

	workers := []WorkerNode{}
	if len(*workerNodeAddrParams) == 0 {
		workers = append(workers, WorkerNode{Address: "localhost:8080"})
	} else {
		for _, w := range strings.Split(*workerNodeAddrParams, ",") {
			workers = append(workers, WorkerNode{Address: w})
		}
	}

	WorkScheduler = NewRandomScheduler(workers)

	http.HandleFunc("/runWorkflow", handleWorkFlow)
	ipAddress := util.GetHostIPAddress()

	http.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "Hello from master::: %s!", ipAddress)
	})

	//Listen on all network devices
	addr := ":8081"
	fmt.Println("Starting master node at:::" + addr)
	http.ListenAndServe(addr, nil)
}
