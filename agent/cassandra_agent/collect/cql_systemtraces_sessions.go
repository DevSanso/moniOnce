package collect

import (
	"agent_common/pkg/applnew/logger"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
	"context"
)
func CollectCQLSystemTracesSessions(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error){
	if err := ctl.ConnectCQL(); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	ctl.Query(ctx, _CqlSystemTracesQuery, func(scanFn func(...any) error) (any, error) {



		return nil, nil
	})

	return nil, nil

}