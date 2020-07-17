package workflow

import (
	"io/ioutil"
	"path"
	"workhorse/pkg/api"

	"gopkg.in/yaml.v2"
)

func ReadWorkflow(workflowFile string) (*api.WorkflowTransferObject, error) {
	workflow, err := readWorkflowYamlFile(workflowFile)
	if err != nil {
		return nil, err
	}
	return convertToTransferObject(workflowFile, workflow)
}

func convertToTransferObject(workflowFile string, workflow *api.Workflow) (*api.WorkflowTransferObject, error) {

	basePath := path.Dir(workflowFile)
	workFlowTransferObject := &api.WorkflowTransferObject{}

	for _, job := range workflow.Jobs {
		//Read contents of script
		script, err := ioutil.ReadFile(path.Join(basePath, job.Script))
		if err != nil {
			return nil, err
		}

		//Convert the information to Trasnfer object
		workFlowTransferObject.Jobs = append(workFlowTransferObject.Jobs, api.JobTransferObject{
			Name:           job.Name,
			ScriptContents: script,
			FileName:       job.Script,
			Image:          job.Image,
		})
	}

	return workFlowTransferObject, nil
}

func readWorkflowYamlFile(workflowFile string) (*api.Workflow, error) {
	data, err := ioutil.ReadFile(workflowFile)
	if err != nil {
		return nil, err
	}

	workFlow := &api.Workflow{}
	err = yaml.Unmarshal(data, workFlow)

	if err != nil {
		return nil, err
	}

	return workFlow, nil
}
