CREATE TABLE mcsa_system_logs (
    object_id bigint,
    collect_time timestamp,
    log_time timestamp,
    order_in_millisecond int,
    "level" text,
    logger text,
    "message" text,
    agent_version float,
    PRIMARY KEY ((object_id, agent_version), collect_time, log_time)
) WITH CLUSTERING ORDER BY (collect_time DESC, log_time DESC)
