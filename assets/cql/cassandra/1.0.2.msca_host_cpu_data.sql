CREATE TABLE msca_host_cpu_data (
    object_id bigint,
    collect_time timestamp,
    sys float,
    user float,
    wait float,
    idle float,
    PRIMARY KEY (object_id, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);