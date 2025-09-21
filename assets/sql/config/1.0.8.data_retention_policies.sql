create table data_retention_policies (
    object_and_sum_id bigint,
    table_name varchar(1024),
    retention_range_day int
);