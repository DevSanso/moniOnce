package collect

import (
	"agent_common/pkg/applnew/logger"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
	"context"
)

func CollectSystemLocalhandle(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error){
	return nil, nil
}