package scheduler

import (
	"encoding/json"
	"fmt"
	eventlister "workhorse/pkg/server/eventlistener"
)

type SchedulerListener struct {
	EventChan chan []interface{}
}

func NewSchedulerListener() *SchedulerListener {
	return &SchedulerListener{
		EventChan: make(chan []interface{}),
	}
}

func (listener *SchedulerListener) Receive(event string) {
	eventInfo := eventlister.EventObject{}
	json.Unmarshal([]byte(event), &eventInfo)

	if eventInfo.Table == "build" {
		fmt.Println(eventInfo.Data)
		listener.EventChan <- []interface{}{eventInfo.Data["id"], eventInfo.Data["created_ts"], eventInfo.Data["status"]}
	}

}
