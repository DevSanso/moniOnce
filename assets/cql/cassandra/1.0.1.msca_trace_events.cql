CREATE TABLE msca_trace_events (
    object_id bigint,
    collect_time timestamp,
    session_id uuid,
    event_id timeuuid,
    activity text,
    source inet,
    source_elapsed int,
    source_port int,
    thread text,
    agent_version float,
    PRIMARY KEY ((object_id, agent_version), collect_time, session_id, event_id)
) WITH CLUSTERING ORDER BY (collect_time DESC);