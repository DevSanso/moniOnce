package collect

import (
	apptype "agent_common/pkg/applnew/types"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
)

var CollectMapping = map[string]apptype.CollectFn[types.PushData, *cassandra.CassandraConn] {
	"collect.nodetool.tpstats" : CollectNodeToolTpStats,
	"collect.nodetool.info"    : CollectNodeToolInfo,
	"collect.cql.trace_session" : CollectCQLSystemTracesSessions,
	"collect.cql.running_query"    : CollectCQLSystemViewRunningQuery,
	"collect.agent.host_cpu" : CollectAgentHostCpu,
	"collect.agent.host_mem" : CollectAgentHostMem,
}

const (
	_CqlSystemTracesQuery = " SELECT session_id, client, command, coordinator, coordinator_port, " + 
		" duration, parameters, request, started_at " + 
		" FROM system_traces.sessions " +
		" WHERE started_at > toTimestamp(now()) - 1000 * 60 " +
		" ALLOW FILTERING "
	
	_CqlSystemViewQueriesQuery = " select thread_id, queued_micros, running_micros, task from system_views.queries "
)
