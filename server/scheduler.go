package main

import "workhorse/util"

type Scheduler interface {
	GetNext() WorkerNode
}

type WorkerNode struct {
	Address string
}

type RandomScheduler struct {
	Nodes []WorkerNode
}

func NewRandomScheduler(workNodes []WorkerNode) *RandomScheduler {
	return &RandomScheduler{Nodes: workNodes}
}

func (random *RandomScheduler) GetNext() WorkerNode {
	totalNodes := len(random.Nodes)
	idx := util.RandomBetween(0, totalNodes)
	return random.Nodes[idx]
}
