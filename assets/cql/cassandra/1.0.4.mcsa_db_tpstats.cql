create table mcsa_db_tpstats (
    object_id bigint,
    collect_time timestamp,
    pool_name text,
    active bigint,
    pending bigint,
    completed bigint,
    blocked bigint

    primary key(object_id, pool_name, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);