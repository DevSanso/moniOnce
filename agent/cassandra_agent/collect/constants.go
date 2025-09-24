package collect

import (
	apptype "agent_common/pkg/applnew/types"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
)

var CollectMapping = map[string]apptype.CollectFn[types.PushData, *cassandra.CassandraConn] {
	"system_local" : CollectSystemLocalhandle,
}