CREATE TABLE msca_jvm_mem_data (
    object_id bigint,
    collect_time timestamp,
    heap_init bigint,
    heap_used bigint,
    heap_committed bigint,
    heap_max bigint,
    nonheap_init bigint,
    nonheap_used bigint,
    nonheap_committed bigint,
    nonheap_max bigint,
    PRIMARY KEY (object_id, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);