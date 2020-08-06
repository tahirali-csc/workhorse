package main

import (
	"flag"
	"workhorse/pkg/client/workflow"
)

func main() {
	//Read command line arguments
	masterNodeAddressParam := flag.String("master-node-address", "localhost", "Address of master node")
	workflowFileParam := flag.String("workflow-file", "", "Workflow file")
	flag.Parse()

	masterNodeAddress := *masterNodeAddressParam
	workflowPath := *workflowFileParam

	//Read workflow
	reader := workflow.Reader{}
	workflowTransferObj, err := reader.ReadWorkflow(workflowPath)
	if err != nil {
		panic(err)
	}

	//Execute workflow
	executor := workflow.Executor{}
	executor.Execute(masterNodeAddress, workflowTransferObj)
}
