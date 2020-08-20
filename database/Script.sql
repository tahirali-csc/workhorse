drop table build


create table build
(
	id INT GENERATED ALWAYS AS identity,
	status varchar(30),
	project_id int,
	created_ts timestamp,
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
	created_ts timestamp,
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


	CREATE OR REPLACE FUNCTION public.notify_event()
	 RETURNS trigger
	 LANGUAGE plpgsql
	AS $function$

	    DECLARE
	        data json;
	        notification json;

	    BEGIN

	        -- Convert the old or new row to JSON, based on the kind of action.
	        -- Action = DELETE?             -> OLD row
	        -- Action = INSERT or UPDATE?   -> NEW row
	        IF (TG_OP = 'DELETE') THEN
	            data = row_to_json(OLD);
	        ELSE
	            data = row_to_json(NEW);
	        END IF;

	        -- Contruct the notification as a JSON string.
	        notification = json_build_object(
	                          'table',TG_TABLE_NAME,
	                          'action', TG_OP,
	                          'data', data);


	        -- Execute pg_notify(channel, notification)
	        PERFORM pg_notify('events',notification::text);

	        -- Result is ignored since this is an AFTER trigger
	        RETURN NULL;
	    END;
	$function$
	;



CREATE TRIGGER build_notify_event
AFTER INSERT OR UPDATE OR DELETE ON build
    FOR EACH ROW EXECUTE PROCEDURE notify_event()
