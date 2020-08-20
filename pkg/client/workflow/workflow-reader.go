package workflow

import (
	"io/ioutil"
	"workhorse/pkg/api"

	inputApi "workhorse/pkg/client/api"

	"gopkg.in/yaml.v2"
)

type Reader struct {
}

func (reader *Reader) ReadWorkflow(workflowFile string) (*api.Workflow, error) {
	workflow, err := reader.readWorkflowYamlFile(workflowFile)
	if err != nil {
		return nil, err
	}
	return reader.mapToWorkflowObj(workflowFile, workflow)
}

func (reader *Reader) mapToWorkflowObj(workflowFile string, workflowInput *inputApi.WorkflowInput) (*api.Workflow, error) {

	// basePath := path.Dir(workflowFile)
	workflow := &api.Workflow{}

	for _, job := range workflowInput.Jobs {
		//Read contents of script
		// _, err := ioutil.ReadFile(path.Join(basePath, job.Script))
		// if err != nil {
		// 	return nil, err
		// }

		//Convert the information to Workflow object
		workflow.Jobs = append(workflow.Jobs, api.WorkflowJob{
			Name: job.Name,
			// ScriptContents: script,
			FileName: job.Script,
			Image:    job.Image,
		})
	}

	return workflow, nil
}

func (reader *Reader) readWorkflowYamlFile(workflowFile string) (*inputApi.WorkflowInput, error) {
	data, err := ioutil.ReadFile(workflowFile)
	if err != nil {
		return nil, err
	}

	workFlow := &inputApi.WorkflowInput{}
	err = yaml.Unmarshal(data, workFlow)

	if err != nil {
		return nil, err
	}

	return workFlow, nil
}
