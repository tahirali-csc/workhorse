package server

import (
	"log"
	"net/url"
	"sync"
	"workhorse/pkg/api"
	"workhorse/pkg/db"
	"workhorse/pkg/server/buildlogs"
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

	// dataChan := make(chan []byte)

	var wg sync.WaitGroup
	wg.Add(1)

	// const baseDir = "/Users/tahir/workspace/workhorse-logs/container-logs"

	go func() {

		wtObj := util.ConvertToWorkflowObject(msg)
		// for _, job := range wtObj.Jobs {
		// 	sendJobToWorkerNodeSync(job, scheduler.GetNext(), dataChan)
		// }

		// buildId := db.CreateBuild("Started")

		// for _, job := range wtObj.Jobs {
		// 	folderName := uuid.New()
		// 	jobPath := path.Join(baseDir, "test-app", folderName.String())
		// 	os.MkdirAll(jobPath, 0755)
		// 	file, _ := os.Create(path.Join(jobPath, "logs.txt"))

		// 	jId := db.CreateBuildJob(buildId, job.Name, "Started", file.Name())
		// 	log.Println("DB Job ID::", jId)
		// 	sendJobToWorkerNodeSync(job, scheduler.GetNext(), file)
		// 	db.UpdateBuildJob(jId, "Finished")
		// }

		bid, bbj := db.CreateBuildStructure(wtObj.Jobs)

		for i, job := range wtObj.Jobs {
			sendJobToWorkerNodeSync(job, scheduler.GetNext(), bbj[i].FileLogs)
			db.UpdateBuildJob(bbj[i].JobID, "Finished")
		}

		wg.Done()
		db.UpdateBuild(bid, "End")
		clientConn.Close()
	}()

	// file, err := os.Create("/Users/tahir/workspace/workhorse-logs/log1.txt")
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println("Temp file::", file.Name())
	// defer file.Close()

	wg.Wait()

	//Stream the response and send to client
	// for msg := range dataChan {
	// 	clientConn.WriteMessage(websocket.BinaryMessage, msg)
	// 	file.WriteString(string(msg))
	// }

	log.Println("Finished the worflow")
}

func sendJobToWorkerNodeSync(job api.WorkflowJob, worker WorkerNode, logFile *buildlogs.FileLogs) {

	log.Println("Sending the job at " + worker.Address)
	u := url.URL{Scheme: "ws", Host: worker.Address, Path: "/runJob"}

	workerNodeConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	jobByteArray, _ := util.ConvertToByteArray(job)
	//Convert the job object to byte array
	workerNodeConn.WriteMessage(websocket.BinaryMessage, jobByteArray)

	for {
		//Read the response from worker node
		msgType, msg, err := workerNodeConn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.CloseMessage {
			break
		} else {
			// logFile.WriteString(string(msg))
			logFile.Write(msg)
		}
	}

	defer func() {
		log.Println("Finished executing job")
		workerNodeConn.Close()
	}()
}

// func sendJobToWorkerNodeSync(job api.WorkflowJob, worker WorkerNode, dataChan chan []byte) {

// 	log.Println("Sending the job at " + worker.Address)
// 	u := url.URL{Scheme: "ws", Host: worker.Address, Path: "/runJob"}

// 	workerNodeConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	//Convert the job object to byte array
// 	workerNodeConn.WriteMessage(websocket.BinaryMessage, util.ConvertToByteArray(job))

// 	for {
// 		//Read the response from worker node
// 		msgType, msg, err := workerNodeConn.ReadMessage()
// 		if err != nil {
// 			break
// 		}

// 		if msgType == websocket.CloseMessage {
// 			break
// 		} else {
// 			dataChan <- msg
// 		}
// 	}

// 	defer func() {
// 		log.Println("Finished executing job")
// 		workerNodeConn.Close()
// 	}()
// }
