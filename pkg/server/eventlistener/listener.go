package eventlister

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type EventObject struct {
	Table  string                 `json: table`
	Action string                 `json: action`
	Data   map[string]interface{} `json:data`
}

type EventListener interface {
	Receive(string)
}

type DBEventsListener struct {
	eventListener []EventListener
}

func (dbListener *DBEventsListener) AddListener(newListener EventListener) {
	dbListener.eventListener = append(dbListener.eventListener, newListener)
}

func (dbListener *DBEventsListener) waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			// fmt.Println("Received data from channel [", n.Channel, "] :")
			// Prepare notification payload for pretty print
			// var prettyJSON bytes.Buffer
			// err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			// if err != nil {
			// 	fmt.Println("Error processing JSON: ", err)
			// 	return
			// }
			// fmt.Println(string(prettyJSON.Bytes()))

			for _, listener := range dbListener.eventListener {
				listener.Receive(n.Extra)
			}
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}

func (dbListener *DBEventsListener) StartListener() {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	// fmt.Println("Start monitoring PostgreSQL...")
	for {
		dbListener.waitForNotification(listener)
	}
}
