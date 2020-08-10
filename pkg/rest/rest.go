package rest

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"workhorse/pkg/api"

	_ "github.com/lib/pq"
)

func GetBuildJobs(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	buildId := request.URL.Query().Get("buildId")
	b, _ := strconv.Atoi(buildId)
	json, _ := json.Marshal(getBuildJobList(b))
	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
}

func GetProjectListHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	json, _ := json.Marshal(getProjectList())
	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
}

func GetBuildLogs(w http.ResponseWriter, request *http.Request) {
	log.Println("Reading build logs...")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	// w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Transfer-Encoding", "chunked")
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	buildId, _ := strconv.Atoi(request.URL.Query().Get("buildId"))
	getBuildLogs(buildId, w, f)

	w.Write(formatSSE("close", ""))
	f.Flush()
}

func GetProjectBuilds(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Access-Control-Allow-Origin", "*")

	pid := request.URL.Query().Get("projectId")
	fmt.Println(pid)

	prjId, _ := strconv.Atoi(pid)
	json, _ := json.Marshal(getProjectBuild(prjId))
	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
}

func getProjectList() []api.Project {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	selectStmt := `
	SELECT * FROM project
	`

	rows, _ := db.Query(selectStmt)
	defer rows.Close()

	plist := []api.Project{}
	for rows.Next() {
		var id int
		var name string

		rows.Scan(&id, &name)

		plist = append(plist, api.Project{
			ID:   id,
			Name: name,
		})
		// fmt.Println(id, name)
	}

	return plist
}

func getProjectBuild(projectId int) []api.ProjectBuild {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	selectStmt := `
	SELECT id, status, start_ts, end_ts FROM build WHERE project_id=$1
	ORDER BY start_ts DESC 
	`

	rows, _ := db.Query(selectStmt, projectId)
	defer rows.Close()

	plist := []api.ProjectBuild{}
	for rows.Next() {
		var id int
		var status string
		var startTs, endTs time.Time

		rows.Scan(&id, &status, &startTs, &endTs)

		plist = append(plist, api.ProjectBuild{
			ID:      id,
			StartTs: startTs,
			EndTs:   endTs,
			Status:  status,
		})

	}

	return plist
}

func getBuildLogs(buildId int, response http.ResponseWriter, f http.Flusher) {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	selectStmt := `
	SELECT id, status, build_log_file FROM build_jobs WHERE build_id=$1
	`

	rows, _ := db.Query(selectStmt, buildId)
	defer rows.Close()

	for rows.Next() {
		var status, build_log_file string
		var bj_id int

		err := rows.Scan(&bj_id, &status, &build_log_file)
		if err != nil {
			log.Fatal(err)
		}

		file, _ := os.Open(build_log_file)
		r := bufio.NewReader(file)

		statChann := make(chan bool)
		ticker := time.NewTicker(time.Second * 10)
		go func() {
			checkStatus := func() bool {
				checkStatus := `
						SELECT status FROM build_jobs WHERE id=$1
					`
				var _status string
				row := db.QueryRow(checkStatus, bj_id)
				row.Scan(&_status)

				if _status == "Finished" {
					return true
				}

				return false
			}

			if checkStatus() {
				statChann <- true
				return
			}

			for {
				select {
				case <-ticker.C:
					if checkStatus() {
						statChann <- true
						return
					}
					// checkStatus := `
					// 	SELECT status FROM build_jobs WHERE id=$1
					// `
					// var _status string
					// row := db.QueryRow(checkStatus, bj_id)
					// row.Scan(&_status)

					// if _status == "Finished" {
					// 	statChann <- true
					// 	return
					// }
				}
			}
		}()

	loop:
		for {
			by, err := r.ReadBytes('\n')
			if err == nil {
				line := string(by)
				// response.Write(formatSSE("message", line))
				// response.Write([]byte(line))
				log.Println("Sending::", line)
				fmt.Fprintf(response, "id: %d\n", bj_id)
				fmt.Fprintf(response, "data: %s\n\n", line)
				f.Flush()
			}

			select {
			case <-statChann:
				break loop
			default:
			}
		}

		for {
			by, err := r.ReadBytes('\n')
			if err != nil {
				log.Println("Have read logs")
				break
			}

			line := string(by)
			// response.Write(formatSSE("message", line))
			// response.Write(
			fmt.Fprintf(response, "id: %d\n", bj_id)
			fmt.Fprintf(response, "data: %s\n\n", line)

			f.Flush()
		}

		log.Println("Going to read next job")
	}

}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	// dataLines := strings.Split(data, "\n")
	// for _, line := range dataLines {
	// 	eventPayload = eventPayload + "data: " + line + "\n"
	// }
	eventPayload = eventPayload + "data: " + data + "\n\n"
	// eventPayload = eventPayload + "data: " + data + "\n"
	return []byte(eventPayload + "\n\n")
}

func getBuildJobList(buildId int) []api.BuildJobInfo {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	selectStmt := `
	SELECT id,job_name FROM build_jobs
	WHERE build_id=$1
	ORDER BY id ASC
	`

	rows, _ := db.Query(selectStmt, buildId)
	defer rows.Close()

	plist := []api.BuildJobInfo{}
	for rows.Next() {
		var id int
		var name string

		rows.Scan(&id, &name)

		plist = append(plist, api.BuildJobInfo{
			Id:      id,
			JobName: name,
		})
	}

	return plist
}
