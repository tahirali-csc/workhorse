package main

import (
	"flag"
	"workhorse/pkg/client/workflow"
	"workhorse/pkg/util"
)

func main() {
	//Read command line arguments
	masterNodeAddressParam := flag.String("master-node-address", "localhost", "Address of master node")
	workflowFileParam := flag.String("workflow-file", "", "Workflow file")
	flag.Parse()

	masterNodeAddress := *masterNodeAddressParam
	workflowPath := *workflowFileParam

	//Read workflow
	workflowTransferObj, err := workflow.ReadWorkflow(workflowPath)
	if err != nil {
		panic(err)
	}

	workflowTransferObjBytes := util.ConvertToByteArray(workflowTransferObj)
	workflow.SendWorkFlow(masterNodeAddress, workflowTransferObjBytes)

}
