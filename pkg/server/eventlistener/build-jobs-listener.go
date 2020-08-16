package eventlister

import (
	"encoding/json"
	"fmt"
	"sync"
)

type BuildJobsEventListener struct {
	DataChannel   chan BuildEventObject
	Cache         BuildJobCache
	EventChannels map[int](chan []string)
	sync          sync.Mutex
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

func (listener *BuildJobsEventListener) Add(id int) chan []string {
	defer listener.sync.Unlock()
	listener.sync.Lock()

	_, ok := listener.EventChannels[id]
	if !ok {
		listener.EventChannels[id] = make(chan []string)
	}
	return listener.EventChannels[id]
}

//TODO: Will review why not pointer receiver???
func (listerner *BuildJobsEventListener) Receive(event string) {
	eventInfo := EventObject{}
	json.Unmarshal([]byte(event), &eventInfo)
	// log.Println("Event Receiver::", eventInfo)

	idVal := int((eventInfo.Data["id"]).(float64))

	if eventInfo.Table == "build_jobs" {
		listerner.Cache.Add(idVal, []string{eventInfo.Data["status"].(string), eventInfo.Data["build_log_file"].(string)})

		// log.Println("Sending data to channel---", idVal)
		defer listerner.sync.Unlock()
		listerner.sync.Lock()

		_, ok := listerner.EventChannels[idVal]
		if !ok {
			listerner.EventChannels[idVal] = make(chan []string)
		}

		// fmt.Println(listerner.EventChannels[idVal])
		msg := []string{eventInfo.Data["status"].(string), eventInfo.Data["build_log_file"].(string)}
		select {
		case listerner.EventChannels[idVal] <- msg:
		default:
		}

	}
}
