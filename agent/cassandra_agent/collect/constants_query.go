package collect

const (
	_CqlSystemTracesQuery = " SELECT session_id, client, command, coordinator, coordinator_port, " +
		" duration, parameters, request, started_at " +
		" FROM system_traces.sessions " +
		" WHERE started_at > toTimestamp(now()) - 1000 * ? " +
		" ALLOW FILTERING "

	_CqlSystemViewQueriesQuery = " select thread_id, queued_micros, running_micros, task from system_views.queries "
	_CqlSystemViewClientsQuery = " select address, driver_name, connection_stage, hostname, request_count, keyspace_name, username from system_views.clients "
	_CqlSystemViewSystemLogsQuery = "  select timestamp, level, logger, message from system_logs where timestamp  > toTimestamp(now()) - 1000 * ?  "
)
