package eventlister

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type BuildJobsEventListener struct {
	DataChannel chan BuildEventObject
	Cache       BuildJobCache
}

type BuildEventObject struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	File   string `json:"build_log_file"`
}

type BuildJobCache struct {
	// data map[int][]string
	s sync.Map
}

func NewBuildJobCache() BuildJobCache {
	return BuildJobCache{
		// data: make(map[int][]string),
	}
}

func (c *BuildJobCache) Add(key int, values []string) {
	// c.s.Lock()
	// c.data[key] = values
	fmt.Println("Cache Entry::", key, "---->", values)
	c.s.Store(key, values)
	// c.s.Unlock()
}

func (c *BuildJobCache) Get(key int) []string {
	// c.s.Lock()
	// defer c.s.Unlock()

	// return c.data[key]
	val, _ := c.s.Load(key)
	return val.([]string)
}

func (c *BuildJobCache) Contains(key int) bool {
	// c.s.Lock()
	// defer c.s.Unlock()

	// _, ok := c.data[key]
	// return ok

	_, ok := c.s.Load(key)
	return ok
}

//TODO: Will review why not pointer receiver???
func (listerner *BuildJobsEventListener) Receive(event string) {
	eventInfo := EventObject{}
	json.Unmarshal([]byte(event), &eventInfo)
	log.Println("Event Receiver::", eventInfo)

	idVal := int((eventInfo.Data["id"]).(float64))

	if eventInfo.Table == "build_jobs" {
		listerner.Cache.Add(idVal, []string{eventInfo.Data["status"].(string), eventInfo.Data["build_log_file"].(string)})
		// listerner.DataChannel <- BuildEventObject{
		// 	Id:     idVal,
		// 	Status: eventInfo.Data["status"].(string),
		// 	File:   eventInfo.Data["build_log_file"].(string),
		// }
	}
}
