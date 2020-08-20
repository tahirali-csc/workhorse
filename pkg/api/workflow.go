package api

import "time"

type Workflow struct {
	Jobs []WorkflowJob
}

type WorkflowJob struct {
	ID       int
	Name     string
	FileName string
	// ScriptContents []byte
	Image string
}

type NodeStats struct {
	*CPUStats
	*MemoryStats
}

type CPUStats struct {
	CPULoad float64 `json:"cpuLoad"`
}

type MemoryStats struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
}

type NodeInfo struct {
	IP string
	MemoryStats
	LastUpdated time.Time
}

type ServerConfig struct {
	Address string
	Port    uint
}

type Project struct {
	ID   int
	Name string
}

type ProjectBuild struct {
	ID      int
	Status  string
	StartTs time.Time
	EndTs   time.Time
}

type BuildJobInfo struct {
	Id      int    `json:"id"`
	JobName string `json:"name"`
}
