package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"workhorse/pkg/api"
	"workhorse/pkg/rest"
	"workhorse/pkg/server"
	as1 "workhorse/pkg/server/api"
	"workhorse/pkg/util"

	eventlister "workhorse/pkg/server/eventlistener"

	scheuduler "workhorse/pkg/scheuduler"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var WorkScheduler server.Scheduler = nil
var sm = server.StatsManager{}
var serverConfig = as1.ServerConfig{}

func handleWorkFlow(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		panic(err)
	}

	server.RunWorkFlowSync(conn, WorkScheduler, serverConfig)

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

		// log.Println(fmt.Sprintf("[%s] Memory : Free=%f, Total=%f, Used=%f --CPU Load=%f",
		// 	senderIP,
		// 	mi.MemoryStats.Free, mi.MemoryStats.Total, mi.MemoryStats.Used,
		// 	mi.CPUStats.CPULoad))
	}
}

func handleReadLogs(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	data, _ := ioutil.ReadFile("/Users/tahir/workspace/worklogs/dat1")
	fmt.Print(string(data))
	fmt.Fprint(response, string(data))
}

func registerAPIEndPoints() {
	http.HandleFunc("/projectList", rest.GetProjectListHandler)
	http.HandleFunc("/projectBuilds", rest.GetProjectBuilds)
	http.HandleFunc("/buildLogs", func(r http.ResponseWriter, rq *http.Request) {
		rest.GetBuildLogs(buildJobListener, r, rq)
	})

	http.HandleFunc("/buildJobs", rest.GetBuildJobs)

	http.HandleFunc("/tempFile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		contents, _ := ioutil.ReadFile("/Users/tahir/workspace/workhorse-logs/container-logs/9bc01d59-5d73-4a5a-9d97-6080cb306fde/logs.txt")
		fmt.Fprintf(w, string(contents))
	})
}

var buildJobListener = &eventlister.BuildJobsEventListener{
	DataChannel:   make(chan eventlister.BuildEventObject),
	Cache:         eventlister.NewBuildJobCache(),
	EventChannels: make(map[int]chan []string),
}

func main() {
	//Read command line arguments
	scheduleParam := flag.String("scheduler", "random", "Worker node schdeduler")
	containerLogsFolderParam := flag.String("containerLogsFolder", "", "Container Logs folder")
	flag.Parse()

	lister := &server.WorkerNodeLister{StatsManager: &sm}

	if *scheduleParam == "random" {
		WorkScheduler = server.NewRandomScheduler(lister)
	} else {
		WorkScheduler = server.NewMemoryBasedScheduler(lister)
	}

	serverConfig.ContainerLogsFolder = *containerLogsFolderParam

	dbEventsListener := eventlister.DBEventsListener{}
	dbEventsListener.AddListener(buildJobListener)

	go dbEventsListener.StartListener()

	jobScheduler := scheuduler.JobScheduler{}
	go jobScheduler.Start(&dbEventsListener, serverConfig, WorkScheduler)

	http.HandleFunc("/runWorkflow", handleWorkFlow)
	http.HandleFunc("/nodestats", handleNodeStats)
	ipAddress := util.GetHostIPAddress()

	http.HandleFunc("/ping", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "Hello from master::: %s!", ipAddress)
	})

	http.HandleFunc("/read", handleReadLogs)
	registerAPIEndPoints()

	//Listen on all network devices
	addr := ":8081"
	log.Println("Starting master node at:::" + addr)
	http.ListenAndServe(addr, nil)
}
