package scheduler

import (
	"database/sql"
	"log"
	"time"
	"workhorse/pkg/api"
	"workhorse/pkg/server"
	"workhorse/pkg/server/buildlogs"
	eventlister "workhorse/pkg/server/eventlistener"
	"workhorse/pkg/util"

	as1 "workhorse/pkg/server/api"
)

type JobScheduler struct {
}

func (sch *JobScheduler) Start(dbListener *eventlister.DBEventsListener, config as1.ServerConfig, nodeScheduler server.Scheduler) {
	listener := NewSchedulerListener()
	dbListener.AddListener(listener)
	wqueue := NewWorkQueue()

	go func() {

		allowedJobs := 1
		sm := make(chan int, allowedJobs)
		for {

			if wqueue.Len() > 0 {
				sm <- 1

				id := wqueue.Remove()

				bid := id.([]interface{})[0].(int)
				go func() {
					sch.updateBuildStart(bid)

					buildJobs := sch.getBuildJobDetails(bid)

					for _, job := range buildJobs {
						continerLogsWriter := buildlogs.NewContainerLogsWriter(config)
						logLocation := continerLogsWriter.GetLocation()
						updateBuildJobStatusAndLogLocation(job.ID, "Started", logLocation)
						server.SendJobToWorkerNodeSync(job, nodeScheduler.GetNext(), continerLogsWriter)
						updateBuildJob(job.ID, "Finished")
					}

					sch.updateBuildFinished(bid)
					<-sm
				}()

				// fmt.Println("Queu:::", bid)

				// running++

				// db.UpdateBuildJobStatusAndLogLocation(bbj[i].JobID, "Started", logLocation)
				// server.SendJobToWorkerNodeSync(job, scheduler.GetNext(), continerLogsWriter)
				// db.UpdateBuildJob(bbj[i].JobID, "Finished")

			}
			// }
		}
	}()

	go func() {
		for {
			select {
			case data := <-listener.EventChan:
				id := int(data[0].(float64))
				createdTs := data[1].(string)
				status := data[2].(string)

				if status == "Pending" {
					tm := util.ParseDBStringTime(createdTs)
					log.Println("Scheduler::", id, tm)

					wqueue.Add([]interface{}{id})
				}
			}
		}
	}()
	// buildIds := sch.getPendingBuilds()
}

//
// func (sch *Scheduler) startQueueListener() {
// 	delay, _ := time.ParseDuration("30s")
// 	ticker := time.NewTicker(delay)
//
// 	select {
// 	case <-ticker.C:
// 		buildIds := getPendingBuilds()
// 	}
// }

func (sch *JobScheduler) updateBuildStart(buildId int) {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	updateStmt := `update build
	set status=$1,
		start_ts=$2
	where id=$3`

	_, err = db.Exec(updateStmt, "Started", time.Now(), buildId)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
}

func (sch *JobScheduler) updateBuildFinished(buildId int) {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	updateStmt := `update build
	set status=$1,
		end_ts=$2
	where id=$3`

	_, err = db.Exec(updateStmt, "Finished", time.Now(), buildId)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
}

func (sch *JobScheduler) getPendingBuilds() []int {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	selectStmt := `select b.id
	from build b
	where b.status = 'Pending'
	order by b.created_ts asc`

	rows, _ := db.Query(selectStmt)
	defer rows.Close()

	var buildId []int
	for rows.Next() {
		var id int

		rows.Scan(&id)

		buildId = append(buildId, id)
	}

	return buildId
}

func (sch *JobScheduler) getBuildJobDetails(buildId int) []api.WorkflowJob {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	selectStm := `
	select bj.id, bj.job_name
	from build_jobs bj
	inner join build b on b.id = bj.build_id 
	where b.id = $1
	`

	rows, _ := db.Query(selectStm, buildId)
	defer rows.Close()

	var workflowJobs []api.WorkflowJob
	for rows.Next() {
		var id int
		var jobName string

		rows.Scan(&id, &jobName)

		workflowJobs = append(workflowJobs, api.WorkflowJob{
			ID:       id,
			Name:     jobName + ".sh",
			FileName: jobName + ".sh",
			Image:    "alpine",
		})
	}

	return workflowJobs
}

func updateBuildJob(buildJobId int, status string) int {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	updateStmt := `
	UPDATE build_jobs 
	SET
	status=$1,
	end_ts=$2
	where id=$3
	`

	id := -1
	_, err = db.Exec(updateStmt, status, time.Now(), buildJobId)
	log.Println("ID:::", id)
	if err != nil {
		log.Println(err)
	}

	return id
}

func updateBuildJobStatusAndLogLocation(buildJobId int, status string, file string) int {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	updateStmt := `
	UPDATE build_jobs 
	SET
	start_ts=$1,
	status=$2,
	build_log_file=$3
	where id=$4
	`

	id := -1
	_, err = db.Exec(updateStmt, time.Now(), status, file, buildJobId)
	log.Println("ID:::", id)
	if err != nil {
		log.Println(err)
	}

	return id
}
