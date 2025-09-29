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
	"collect.cql.systemlog": collectCQLSystemViewSystemLogs,
	"collect.agent.host_cpu":    collectAgentHostCpu,
	"collect.agent.host_mem":    collectAgentHostMem,
}
