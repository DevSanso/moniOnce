create table mcsa_trace_session (
    object_id bigint,
    collect_time timestamp,
    session_id text,
    client inet,
    command text,
    coordinator inet,
    coordinator_port int,
    duration int,
    request text,
    started_at timestamp,
    param_1 text,
    param_2 text,

    primary key(object_id, collect_time, session_id)
) WITH CLUSTERING ORDER BY (collect_time DESC);