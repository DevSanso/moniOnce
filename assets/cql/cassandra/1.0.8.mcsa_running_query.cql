CREATE TABLE mcsa_running_query (
    object_id bigint,
    collect_time timestamp,
    thread_id text,
    queue_micro_sec int,
    running_micro_sec int,
    query text,
    agent_version float,
    PRIMARY KEY ((object_id, agent_version), collect_time, queue_micro_sec, running_micro_sec)
) WITH CLUSTERING ORDER BY (collect_time DESC, running_micro_sec DESC, queue_micro_sec DESC)
