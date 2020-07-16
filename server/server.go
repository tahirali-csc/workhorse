package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"workhorse/api"
	"workhorse/pkg/server"
	"workhorse/pkg/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var WorkScheduler server.Scheduler = nil
var sm = server.StatsManager{}

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

func handleNodeStats(_ http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		senderIP := util.GetSenderIP(request)

		mi := api.NodeStats{}
		body, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()

		if err != nil {
			log.Print(err)
			return
		}

		err = json.Unmarshal(body, &mi)
		if err != nil {
			log.Print(err)
			return
		}

		sm.UpdateStats(senderIP, mi)

		log.Println(fmt.Sprintf("[%s] Memory : Free=%f, Total=%f, Used=%f --CPU Load=%f",
			senderIP,
			mi.MemoryStats.Free, mi.MemoryStats.Total, mi.MemoryStats.Used,
			mi.CPUStats.CPULoad))
	}
}

func handleReadLogs(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	data, _ := ioutil.ReadFile("/Users/tahir/workspace/worklogs/dat1")
	fmt.Print(string(data))
	fmt.Fprint(response, string(data))
}

func main() {
	//Read command line arguments
	//workerNodeAddrParams := flag.String("worker-node-address", "", "Comma separted ip address of worker nodes")
	scheduleParam := flag.String("scheduler", "random", "Worker node schdeduler")
	flag.Parse()

	//workers := []server.WorkerNode{}
	//if len(*workerNodeAddrParams) == 0 {
	//	workers = append(workers, server.WorkerNode{Address: "localhost:8080"})
	//} else {
	//	for _, w := range strings.Split(*workerNodeAddrParams, ",") {
	//		workers = append(workers, server.WorkerNode{Address: w})
	//	}
	//}

	lister := &server.WorkerNodeLister{StatsManager: &sm}

	if *scheduleParam == "random" {
		WorkScheduler = server.NewRandomScheduler(lister)
	} else{
		WorkScheduler = server.NewMemoryBasedScheduler(lister)
	}
	//} else if *scheduleParam == "roundroubin" {
	//	WorkScheduler = server.NewRoundRobinScheduler(lister)


	//log.Println("Will use these worker nodes:::", workers)

	http.HandleFunc("/runWorkflow", handleWorkFlow)
	http.HandleFunc("/nodestats", handleNodeStats)
	ipAddress := util.GetHostIPAddress()

	http.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "Hello from master::: %s!", ipAddress)
	})

	http.HandleFunc("/read", handleReadLogs)

	//Listen on all network devices
	addr := ":8081"
	log.Println("Starting master node at:::" + addr)
	http.ListenAndServe(addr, nil)
}
