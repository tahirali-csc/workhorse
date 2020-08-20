package db

import (
	"context"
	"database/sql"
	"log"
	"time"
	"workhorse/pkg/api"

	_ "github.com/lib/pq"
)

type BuildJobDTO struct {
	JobID int
	// FileLogs *buildlogs.FileLogs
}

func CreateBuildStructure(jobs []api.WorkflowJob) (int, []BuildJobDTO) {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	insertStmt := `
	INSERT INTO build (status, project_id, created_ts)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	buildId := -1
	bbj := []BuildJobDTO{}
	err = tx.QueryRow(insertStmt, "Pending", 1, time.Now()).Scan(&buildId)
	for _, j := range jobs {

		// const baseDir = "/Users/tahir/workspace/workhorse-logs/container-logs"
		// fileLogs := buildlogs.NewFileLogs(baseDir)

		insertStmt := `
		INSERT INTO build_jobs (build_id, job_name, status, build_log_file)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`

		id := -1
		err = tx.QueryRow(insertStmt, buildId, j.Name, "Pending", "").Scan(&id)
		log.Println("Job ID:::", id)
		if err != nil {
			log.Println(err)
		}

		bbj = append(bbj, BuildJobDTO{JobID: id})
	}

	tx.Commit()
	return buildId, bbj
}

func CreateBuild(state string) int {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insertStmt := `
	INSERT INTO build (status, project_id, start_ts)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	id := -1
	err = db.QueryRow(insertStmt, "Started", 1, time.Now()).Scan(&id)
	log.Println("ID:::", id)
	if err != nil {
		log.Println(err)
	}

	return id
}

func UpdateBuild(buildId int, state string) {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	updateStmt := `
	UPDATE build
	SET 
	status=$1,
	end_ts=$2
	where id=$3
	`

	_, err = db.Exec(updateStmt, state, time.Now(), buildId)
	if err != nil {
		log.Println(err)
	}
}

func CreateBuildJob(buildId int, jobName, state, buildPath string) int {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insertStmt := `
	INSERT INTO build_jobs (build_id, job_name, status, build_log_file, start_ts)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	id := -1
	err = db.QueryRow(insertStmt, buildId, jobName, "Started", buildPath, time.Now()).Scan(&id)
	log.Println("ID:::", id)
	if err != nil {
		log.Println(err)
	}

	return id
}

func UpdateBuildJob(buildJobId int, status string) int {
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
