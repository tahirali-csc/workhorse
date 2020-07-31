package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func CreateBuild(state string) int {
	var conninfo string = "dbname=postgres user=dev password=dev host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insertStmt := `
	INSERT INTO build (status, start_ts)
	VALUES ($1, $2)
	RETURNING id
	`

	id := -1
	err = db.QueryRow(insertStmt, "Started", time.Now()).Scan(&id)
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
