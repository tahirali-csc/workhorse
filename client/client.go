package main

import (
	"flag"
	"workhorse/pkg/workflow"
	"workhorse/util"
)

func main() {
	//Read command line arguments
	masterNodeIPParam := flag.String("master-node-address", "localhost", "Address of master node")
	workflowFileParam := flag.String("workflow-file", "", "Workflow file")
	flag.Parse()

	masterNodeIP := *masterNodeIPParam
	workflowPath := *workflowFileParam

	// workflowPath = "/home/tahir/workspace/rnd-projects/workhorse/client/sample-workflow/workflow.yaml"

	//Read workflow
	workflowTransferObj, err := workflow.ReadWorkflow(workflowPath)
	if err != nil {
		panic(err)
	}

	workflowTransferObjBytes := util.ConvertToByteArray(workflowTransferObj)
	workflow.SendWorkFlow(masterNodeIP, workflowTransferObjBytes)

}
