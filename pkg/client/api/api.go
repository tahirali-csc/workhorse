package api

type WorkflowInput struct {
	Jobs []JobInput
}

type JobInput struct {
	Name   string
	Script string
	Image  string
}
