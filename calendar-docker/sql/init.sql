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