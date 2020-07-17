package api

import "time"

type WorkflowTransferObject struct {
	Jobs []JobTransferObject
}

type JobTransferObject struct {
	Name           string
	FileName       string
	ScriptContents []byte
	Image          string
}

type Workflow struct {
	Jobs []Job
}

type Job struct {
	Name   string
	Script string
	Image  string
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