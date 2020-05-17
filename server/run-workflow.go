package main

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

func runWorkFlow(clientConn *websocket.Conn) {

	var wg sync.WaitGroup

	for {
		_, _, err := clientConn.ReadMessage()
		if err != nil {
			break
		}

		var lock sync.RWMutex
		for i := 1; i <= 2; i++ {
			wg.Add(1)
			go func() {
				sendJobToWorkerNode(clientConn, &wg, &lock)
			}()
		}
		break
	}

	wg.Wait()
	fmt.Println("Finished the worflow")
}

func sendJobToWorkerNode(clientConn *websocket.Conn, wg *sync.WaitGroup, lock *sync.RWMutex) {
	// const addr = "localhost:8082"
	const addr = "192.168.56.103:8082"
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
			func() {
				lock.Lock()
				clientConn.WriteMessage(websocket.BinaryMessage, msg)
				defer lock.Unlock()
			}()
		}
	}

	defer func() {
		fmt.Println("Finished executing job")
		workerNodeConn.Close()
		wg.Done()
	}()
}
