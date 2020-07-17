package server

import (
	"log"
	"sort"

	"workhorse/pkg/api"
	"workhorse/pkg/util"
)

type Scheduler interface {
	GetNext() WorkerNode
}

type WorkerNode struct {
	Address string
}

type RandomScheduler struct {
	NodeLister *WorkerNodeLister
}

func NewRandomScheduler(lister *WorkerNodeLister) *RandomScheduler {
	return &RandomScheduler{NodeLister: lister}
}

func (random *RandomScheduler) GetNext() WorkerNode {
	nodes := random.NodeLister.getActiveWorkerNodes()
	totalNodes := len(nodes)
	idx := util.RandomBetween(0, totalNodes)
	node := WorkerNode{Address: nodes[idx].IP}
	log.Println("Randomly choosing :: ", node.Address)
	return node
}

//type RoundRobinSchedule struct {
//	Nodes      []WorkerNode
//	currentIdx int
//	lock       sync.RWMutex
//}
//
//func (round *RoundRobinSchedule) GetNext() WorkerNode {
//	totalNodes := len(round.Nodes)
//
//	round.lock.Lock()
//
//	round.currentIdx++
//	if round.currentIdx > (totalNodes - 1) {
//		round.currentIdx = 0
//	}
//
//	defer round.lock.Unlock()
//	return round.Nodes[round.currentIdx]
//}
//
//func NewRoundRobinScheduler(workNodes []WorkerNode) *RoundRobinSchedule {
//	return &RoundRobinSchedule{Nodes: workNodes, currentIdx: -1}
//}

type MemoryScheduler struct {
	NodeLister *WorkerNodeLister
}

func (mss *MemoryScheduler) GetNext() WorkerNode {
	runnableNodes := mss.NodeLister.getActiveWorkerNodes()
	node := WorkerNode{}

	if len(runnableNodes) > 0 {
		sort.Slice(runnableNodes, func(i, j int) bool {
			if runnableNodes[i].Free < runnableNodes[j].Free {
				return false
			}
			return true
		})

		node.Address = runnableNodes[0].IP
	}
	//fmt.Print("Selected Node::", node)
	return node
}

func NewMemoryBasedScheduler(lister *WorkerNodeLister) *MemoryScheduler {
	return &MemoryScheduler{NodeLister: lister}
}

//func ScheduleRun(ms StatsManager) {
//	timer := time.NewTicker(time.Second * 5)
//	for {
//		select {
//		case <-timer.C:
//
//			var runnableNodes []api.NodeInfo
//
//			ms.statsMap.Range(func(key, value interface{}) bool {
//				ni := value.(api.NodeInfo)
//				duration := time.Now().Sub(ni.LastUpdated)
//				if duration.Seconds() <= 5 {
//					runnableNodes = append(runnableNodes, ni)
//				}
//				return true
//			})
//
//			if len(runnableNodes) > 0 {
//				sort.Slice(runnableNodes, func(i, j int) bool {
//					if runnableNodes[i].Free < runnableNodes[j].Free {
//						return false
//					}
//					return true
//				})
//
//				//res := ""
//				//for _, v := range runnableNodes {
//				//	res += fmt.Sprintf("\nIP :: %s Free :: %f", v.IP, v.Free)
//				//}
//				//log.Print(res)
//			}
//		}
//	}
//}

type WorkerNodeLister struct {
	StatsManager *StatsManager
}

func (wl *WorkerNodeLister) getActiveWorkerNodes() []api.NodeInfo {

	var runnableNodes []api.NodeInfo
	wl.StatsManager.statsMap.Range(func(key, value interface{}) bool {
		ni := value.(api.NodeInfo)
		//duration := time.Now().Sub(ni.LastUpdated)
		//if duration.Seconds() <= 5 {
		runnableNodes = append(runnableNodes, ni)
		//}
		return true
	})

	return runnableNodes
}
