drop table build


create table build
(
	id INT GENERATED ALWAYS AS identity,
	status varchar(30),
	project_id int,
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

create table project
(
	id INT generated always as identity,
	name varchar(200)
)

select * from project

insert into project(name)
values('Project1')

insert into project(name)
values('Project2')

truncate table build
truncate table build_jobs 


select * from build
select * from build_jobs



CREATE TRIGGER build_jobs_notify_event
AFTER INSERT OR UPDATE OR DELETE ON build_jobs
    FOR EACH ROW EXECUTE PROCEDURE notify_event()

drop table build
drop table build_jobs 

UPDATE build
	SET end_ts = current
	
	