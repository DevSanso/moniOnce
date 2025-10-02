CREATE TABLE msca_host_mem_data (
    object_id bigint,
    collect_time timestamp,
    total bigint,
    free bigint,
    used bigint,
    agent_version float,
    PRIMARY KEY ((object_id, agent_version), collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);