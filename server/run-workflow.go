package main

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func runWorkFlowSync(clientConn *websocket.Conn) {

	for {
		_, _, err := clientConn.ReadMessage()
		if err != nil {
			break
		}

		for i := 1; i <= 2; i++ {
			sendJobToWorkerNodeSync(clientConn)
		}
		break
	}

	fmt.Println("Finished the worflow")
}

func sendJobToWorkerNodeSync(clientConn *websocket.Conn) {
	const addr = "192.168.56.103:8080"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/runJob"}

	workerNodeConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}

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
