create table mcsa_clients (
    object_id bigint,
    collect_time timestamp,
    "address" text,
    driver_name text,
    connection_stage text,
    hostname text,
    request_count bigint,
    keyspace_name text,
    username text,
    agent_version float,
    primary key((object_id, agent_version), keyspace_name, username, "address", request_count, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC, request_count DESC);
