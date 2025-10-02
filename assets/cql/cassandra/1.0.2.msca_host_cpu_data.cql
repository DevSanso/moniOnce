CREATE TABLE msca_host_cpu_data (
    object_id bigint,
    collect_time timestamp,
    sys float,
    user float,
    wait float,
    idle float,
    agent_version float,
    PRIMARY KEY ((object_id, agent_version), collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);