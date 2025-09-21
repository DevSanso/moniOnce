CREATE TABLE msca_jvm_host_data (
    object_id bigint,
    collect_time timestamp,
    system_cpu float,
    jvm_process_cpu float,
    free_physical_memory bigint,
    total_physical_memory bigint,
    free_swap_space bigint,
    total_swap_space bigint,
    committed_virtual_memory bigint
    PRIMARY KEY (object_id, collect_time)
) WITH CLUSTERING ORDER BY (collect_time DESC);