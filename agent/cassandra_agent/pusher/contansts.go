package pusher


const (
	_InsertHostCpuQuery = " INSERT INTO msca_host_cpu_data(object_id, collect_time, sys, user, wait, idle, agent_version) " + 
		" VALUES (?,?,?,?,?,?,?, ?) USING TTL ? "
	_InsertHostMemoryQuery = " INSERT INTO msca_host_mem_data(object_id, collect_time, total, free, used, agent_version) " +
		" VALUES (?,?,?,?,?, ?) USING TTL ? "
	_InsertCQLClientsQuery = " INSERT INTO msca_clients(object, collect_time, \"address\", driver_name, connection_stage, hostname, request_count, keyspace_name, username, agent_version) " +
	    " VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) USING TTL ? "
	_InsertCQLRunningQueryQuery = " INSERT INTO mcsa_running_query(object_id, collect_time, thread_id, queue_micro_sec, running_micro_sec, query, agent_version) " + 
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) USING TTL ? "
	_Insert
)
