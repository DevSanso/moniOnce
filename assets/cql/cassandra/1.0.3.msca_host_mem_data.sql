CREATE TABLE msca_host_mem_data (
    object_id bigint,
    collect_time timestamp,
    total bigint,
    free bigint,
    used bigint,
    PRIMARY KEY (object_id, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);