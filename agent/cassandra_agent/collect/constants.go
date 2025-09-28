package collect

import (
	apptype "agent_common/pkg/applnew/types"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
)

var CollectMapping = map[string]apptype.CollectFn[types.PushData, *cassandra.CassandraConn]{
	"collect.nodetool.tpstats":  collectNodetoolTpStats,
	"collect.nodetool.info":     collectNodeToolInfo,
	"collect.cql.trace_session": collectCQLSystemTracesSessions,
	"collect.cql.running_query": collectCQLSystemViewRunningQuery,
	"collect.cql.clients" : collectCQLSystemViewClients,
	"collect.agent.host_cpu":    collectAgentHostCpu,
	"collect.agent.host_mem":    collectAgentHostMem,
}

const (
	_CqlSystemTracesQuery = " SELECT session_id, client, command, coordinator, coordinator_port, " +
		" duration, parameters, request, started_at " +
		" FROM system_traces.sessions " +
		" WHERE started_at > toTimestamp(now()) - 1000 * 60 " +
		" ALLOW FILTERING "

	_CqlSystemViewQueriesQuery = " select thread_id, queued_micros, running_micros, task from system_views.queries "
	_CqlSystemViewClientsQuery = " select address, driver_name, connection_stage, hostname, request_count, keyspace_name, username from system_views.clients "
)
