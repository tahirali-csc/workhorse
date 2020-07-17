package server

import (
	"sync"
	"time"
	"workhorse/pkg/api"
)

type StatsManager struct {
	statsMap sync.Map
}

func (m *StatsManager) UpdateStats(ip string, stats api.NodeStats) {

	nodeInfo := api.NodeInfo{
		MemoryStats: *stats.MemoryStats,
		LastUpdated: time.Now(),
		IP:          ip,
	}

	m.statsMap.Store(ip, nodeInfo)
}
