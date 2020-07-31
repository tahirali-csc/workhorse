create table build
(
	id INT GENERATED ALWAYS AS identity,
	status varchar(30),
	start_ts timestamp,
	end_ts timestamp
)

create table build_jobs
(
	id INT GENERATED ALWAYS AS identity,
	build_id	int,
	job_name	varchar(30),
	status varchar(30),
	build_log_file varchar(300),
	start_ts timestamp,
	end_ts timestamp
)


select * from build
select * from build_jobs

drop table build
drop table build_jobs 

UPDATE build
	SET end_ts = current
	
	