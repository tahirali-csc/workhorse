package server

import (
	"log"
	"net/url"
	"workhorse/api"
	"workhorse/pkg/util"

	"github.com/gorilla/websocket"
)

func RunWorkFlowSync(clientConn *websocket.Conn, scheduler Scheduler) {

	//Read the workflow file
	_, msg, err := clientConn.ReadMessage()
	if err != nil {
		log.Fatal(err)
		return
	}

	dataChan := make(chan []byte)

	//Sequentially send the job to worker node
	go func() {
		wtObj := util.ConvertToWorkflowObject(msg)
		for _, job := range wtObj.Jobs {
			sendJobToWorkerNodeSync(job, scheduler.GetNext(), dataChan)
		}

		defer close(dataChan)
	}()

	// Open file using READ & WRITE permission.
	//file, err := os.OpenFile("/Users/tahir/workspace/worklogs/dat1", os.O_RDWR, 0644)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//defer func() {
	//	file.Sync()
	//	file.Close()
	//}()

	//Stream the response and send to client
	for msg := range dataChan {
		clientConn.WriteMessage(websocket.BinaryMessage, msg)
		//file.Write(msg)
	}

	log.Println("Finished the worflow")
}

func sendJobToWorkerNodeSync(job api.JobTransferObject, worker WorkerNode, dataChan chan []byte) {

	log.Println("Sending the job at " + worker.Address)
	u := url.URL{Scheme: "ws", Host: worker.Address, Path: "/runJob"}

	workerNodeConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Convert the job object to byte array
	workerNodeConn.WriteMessage(websocket.BinaryMessage, util.ConvertToByteArray(job))

	for {
		//Read the response from worker node
		msgType, msg, err := workerNodeConn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.CloseMessage {
			break
		} else {
			dataChan <- msg
		}
	}

	defer func() {
		log.Println("Finished executing job")
		workerNodeConn.Close()
	}()

}

func writeFile() {

}
