package server

import (
	"sync"
	"workhorse/pkg/util"
)

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

type RoundRobinSchedule struct {
	Nodes      []WorkerNode
	currentIdx int
	lock       sync.RWMutex
}

func (round *RoundRobinSchedule) GetNext() WorkerNode {
	totalNodes := len(round.Nodes)

	round.lock.Lock()

	round.currentIdx++
	if round.currentIdx > (totalNodes - 1) {
		round.currentIdx = 0
	}

	defer round.lock.Unlock()
	return round.Nodes[round.currentIdx]
}

func NewRoundRobinScheduler(workNodes []WorkerNode) *RoundRobinSchedule {
	return &RoundRobinSchedule{Nodes: workNodes, currentIdx: -1}
}
