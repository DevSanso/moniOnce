package collect

import (
	apptype "agent_common/pkg/applnew/types"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
)

var CollectMapping = map[string]apptype.CollectFn[types.PushData, *cassandra.CassandraConn] {
	"cql_system_local" : CollectCQLSystemLocalhandle,
	"nodetool_tpstats" : CollectNodeToolTpStats,
	"nodetool_info"    : CollectNodeToolInfo,
	"cql_system_traces_sessions" : CollectCQLSystemTracesSessions,
}

const (
	_CqlSystemTracesQuery = " SELECT session_id, client, command, coordinator, coordinator_port, " + 
		" duration, parameters, request, started_at " + 
		" FROM system_traces.sessions " +
		" WHERE started_at > toTimestamp(now()) - 1000 * 60 " +
		" ALLOW FILTERING "
)
