-- +goose Up
CREATE ROLE app_user WITH
	LOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD '';

CREATE DATABASE calendar
    WITH 
    OWNER = app_user
    ENCODING = 'UTF8'
    LC_COLLATE = 'Russian_Russia.1251'
    LC_CTYPE = 'Russian_Russia.1251'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE calendar
    IS 'OTUS Calendar';

CREATE TABLE event
(
    id bigserial PRIMARY KEY,
    title character varying(256) NOT NULL,
    description character varying(1024),
    location character varying(256),
    start_time timestamp(6) with time zone,
    end_time timestamp(6) with time zone,
    notify_before interval,
    user_id bigint,
    calendar_id bigint,
    created timestamp(6) with time zone,
    last_updated timestamp(6) with time zone
);

ALTER TABLE event OWNER to app_user;

-- +goose Down
DROP TABLE event;
