package main

import (
	"flag"
	"workhorse/util"
)

// func main() {
// 	t := time.Now()
// 	fmt.Printf("%d%02d%02d%02d%02d%02d",
// 		t.Year(), t.Month(), t.Day(),
// 		t.Hour(), t.Minute(), t.Second())
// }

func main() {
	//Read command line arguments
	masterNodeIPParam := flag.String("master-node-address", "localhost", "Address of master node")
	workflowFileParam := flag.String("workflow-file", "", "Workflow file")
	flag.Parse()

	masterNodeIP := *masterNodeIPParam
	workflowPath := *workflowFileParam

	// workflowPath = "/home/tahir/workspace/rnd-projects/workhorse/client/sample-workflow/workflow.yaml"

	//Read workflow
	workflowTransferObj, err := readWorkflow(workflowPath)
	if err != nil {
		panic(err)
	}

	workflowTransferObjBytes := util.ConvertToByteArray(workflowTransferObj)
	sendWorkFlow(masterNodeIP, workflowTransferObjBytes)

}
