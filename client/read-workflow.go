package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"workhorse/api"

	"gopkg.in/yaml.v2"
)

const workflowPath = "/home/tahir/workspace/rnd-projects/workhorse/client/sample-workflow"

func readWorkflow() *api.WorkflowTransferObject {
	return convertToTransferObject(readWorkflowYamlFile())
}

func convertToTransferObject(workflow *api.Workflow) *api.WorkflowTransferObject {

	workFlowTransferObject := &api.WorkflowTransferObject{}
	for _, v := range workflow.Jobs {
		data, _ := ioutil.ReadFile(path.Join(workflowPath, v.Script))

		workFlowTransferObject.Jobs = append(workFlowTransferObject.Jobs, api.JobTransferObject{
			Name:           v.Name,
			ScriptContents: data,
		})
	}

	return workFlowTransferObject
}

func readWorkflowYamlFile() *api.Workflow {
	data, _ := ioutil.ReadFile(path.Join(workflowPath, "workflow.yaml"))

	workFlow := &api.Workflow{}
	err := yaml.Unmarshal(data, workFlow)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return workFlow
}
