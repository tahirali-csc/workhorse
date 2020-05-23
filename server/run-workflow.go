package main

import (
	"fmt"
	"net/url"
	"workhorse/api"
	"workhorse/util"

	"github.com/gorilla/websocket"
)

func runWorkFlowSync(clientConn *websocket.Conn, scheduler Scheduler) {

	for {
		_, msg, err := clientConn.ReadMessage()
		if err != nil {
			break
		}

		wtObj := util.ConvertToWorkflowObject(msg)
		for _, job := range wtObj.Jobs {
			sendJobToWorkerNodeSync(clientConn, job, scheduler.GetNext())
		}
		break
	}

	fmt.Println("Finished the worflow")
}

func sendJobToWorkerNodeSync(clientConn *websocket.Conn, job api.JobTransferObject, worker WorkerNode) {
	// const addr = "localhost:8080"
	// const addr = "192.168.56.103:8080"
	fmt.Println("Sending the job at " + worker.Address)
	u := url.URL{Scheme: "ws", Host: worker.Address, Path: "/runJob"}

	workerNodeConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

	workerNodeConn.WriteMessage(websocket.BinaryMessage, util.ConvertToByteArray(job))

	for {
		msgType, msg, err := workerNodeConn.ReadMessage()
		if err != nil {
			break
		}
		if msgType == websocket.CloseMessage {
			break
		} else {
			clientConn.WriteMessage(websocket.BinaryMessage, msg)
		}
	}

	defer func() {
		fmt.Println("Finished executing job")
		workerNodeConn.Close()
	}()
}

//Runs jobs in parallel
// func runWorkFlowAsync(clientConn *websocket.Conn) {

// 	var wg sync.WaitGroup

// 	for {
// 		_, _, err := clientConn.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		var lock sync.RWMutex
// 		for i := 1; i <= 2; i++ {
// 			wg.Add(1)
// 			go func() {
// 				sendJobToWorkerNode(clientConn, &wg, &lock)
// 			}()
// 		}
// 		break
// 	}

// 	wg.Wait()
// 	fmt.Println("Finished the worflow")
// }

// func sendJobToWorkerNode(clientConn *websocket.Conn, wg *sync.WaitGroup, lock *sync.RWMutex) {
// 	const addr = "192.168.56.103:8080"
// 	u := url.URL{Scheme: "ws", Host: addr, Path: "/runJob"}

// 	workerNodeConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for {
// 		msgType, msg, err := workerNodeConn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		if msgType == websocket.CloseMessage {
// 			break
// 		} else {
// 			func() {
// 				lock.Lock()
// 				clientConn.WriteMessage(websocket.BinaryMessage, msg)
// 				defer lock.Unlock()
// 			}()
// 		}
// 	}

// 	defer func() {
// 		fmt.Println("Finished executing job")
// 		workerNodeConn.Close()
// 		wg.Done()
// 	}()
// }
