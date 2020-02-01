CREATE TABLE banner_statistics
(
    banner_id bigint NOT NULL,
    group_id bigint NOT NULL,
    number_of_shows bigint,
    number_of_clicks bigint,
    PRIMARY KEY(banner_id, group_id)
);