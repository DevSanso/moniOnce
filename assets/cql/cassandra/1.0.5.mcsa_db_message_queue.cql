create table mcsa_db_message_queue (
    object_id bigint,
    collect_time timestamp,
    message_type text,
    "drop"      bigint,
    percentile_50 float,
    percentile_95 float,
    percentile_99 float,
    max           float,
    primary key(object_id, message_type, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);