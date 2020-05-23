package api

type WorkflowTransferObject struct {
	Jobs []JobTransferObject
}

type JobTransferObject struct {
	Name           string
	FileName       string
	ScriptContents []byte
}

type Workflow struct {
	Jobs []Job
}

type Job struct {
	Name   string
	Script string
}
