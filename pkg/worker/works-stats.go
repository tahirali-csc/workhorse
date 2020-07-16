package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"log"
	"net/http"
	"time"
	"workhorse/api"
)

func KeepSendingStats(config *api.ServerConfig) {

	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-timer.C:
			stats, err := GetStats()
			if err != nil {
				log.Panic(err)
				return
			}

			statsJson, _ := json.Marshal(stats)

			url := fmt.Sprintf("http://%s:%d/nodestats", config.Address, config.Port)
			_, err = http.Post(url, "application/json", bytes.NewReader(statsJson))
			if err != nil {
				log.Print(err)
			}
		}
	}

}

func toGB(bytes uint64) float64 {
	return float64(bytes) / float64(1024*1024*1024)
}

func GetStats() (*api.NodeStats, error) {
	mem, err := getMemoryStats()
	if err != nil {
		return nil, err
	}

	cpu, err := getCPUStats()
	if err != nil {
		return nil, err
	}

	return &api.NodeStats{
		CPUStats:    cpu,
		MemoryStats: mem,
	}, nil
}

func getMemoryStats() (*api.MemoryStats, error) {
	memory, err := memory.Get()
	if err != nil {
		return nil, err
	}
	return &api.MemoryStats{
		Total: toGB(memory.Total),
		Used:  toGB(memory.Used),
		Free:  toGB(memory.Free),
	}, nil
}

func getCPUStats() (*api.CPUStats, error) {
	before, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(1) * time.Second)

	after, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	total := float64(after.Total - before.Total)
	cpuUser := float64(after.User-before.User) / total * 100
	cpuSystem := float64(after.System-before.System) / total * 100

	return &api.CPUStats{
		CPULoad: cpuUser + cpuSystem,
	}, nil

	// fmt.Printf("cpu user: %f %%\n", float64(after.User-before.User)/total*100)
	// fmt.Printf("cpu system: %f %%\n", float64(after.System-before.System)/total*100)
	// fmt.Printf("cpu idle: %f %%\n", float64(after.Idle-before.Idle)/total*100)
}
