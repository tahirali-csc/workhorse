package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"workhorse/pkg/server"
	"workhorse/pkg/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var WorkScheduler server.Scheduler = nil

func handleWorkFlow(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		panic(err)
	}

	server.RunWorkFlowSync(conn, WorkScheduler)

	defer func() {
		log.Println("Closing socket connection")
		conn.Close()
	}()
}

func main() {
	//Read command line arguments
	workerNodeAddrParams := flag.String("worker-node-address", "", "Comma separted ip address of worker nodes")
	scheduleParam := flag.String("scheduler", "roundrobin", "Worker node schdeduler")
	flag.Parse()

	workers := []server.WorkerNode{}
	if len(*workerNodeAddrParams) == 0 {
		workers = append(workers, server.WorkerNode{Address: "localhost:8080"})
	} else {
		for _, w := range strings.Split(*workerNodeAddrParams, ",") {
			workers = append(workers, server.WorkerNode{Address: w})
		}
	}

	if *scheduleParam == "random" {
		WorkScheduler = server.NewRandomScheduler(workers)
	} else {
		WorkScheduler = server.NewRoundRobinScheduler(workers)
	}

	log.Println("Will use these worker nodes:::", workers)

	http.HandleFunc("/runWorkflow", handleWorkFlow)
	ipAddress := util.GetHostIPAddress()

	http.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "Hello from master::: %s!", ipAddress)
	})

	//Listen on all network devices
	addr := ":8081"
	log.Println("Starting master node at:::" + addr)
	http.ListenAndServe(addr, nil)
}
