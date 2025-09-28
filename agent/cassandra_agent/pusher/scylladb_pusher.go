package pusher

import (
	"agent_common/pkg/applnew/logger"
	apptype "agent_common/pkg/applnew/types"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"context"
	"fmt"
)

type scyllaDbPusher struct {
	pool apptype.CollectConnPool[*cassandra.CassandraConn]
}

func NewScyllaDbPusher(url string, args ...any) (apptype.DataPusher[types.PushData], error) {
	p, err := cassandra.NewCassandraPool(url, args...)
	if err != nil {
		return nil, err
	}
	return &scyllaDbPusher{pool : p}, nil
}

func (sdp *scyllaDbPusher) Push(data *types.PushData, ctx context.Context, log logger.LevelLogger) error {
	conn, connErr := sdp.pool.GetDbConn(ctx)
	if connErr != nil {
		log.Error(connErr.Error())
		return connErr
	}
	
	var pushErr error = nil 

	switch constants.DataID(data.DataId) {
	case constants.AgentHostCpu:
		pushErr = pushHostCpu(ctx, conn, &data.Host.CpuPercent)
	case constants.AgentHostMemory:
		pushErr = pushHostMemory(ctx, conn, &data.Host.Memory)
	case constants.CQLClients:
		pushErr = pushClients(ctx, conn, data.Cql.Clients)
	}

	if pushErr != nil {
		e := fmt.Errorf("insert failed conn_id:%d, data_id:%d, err:%s", data.ConnTypeId, data.DataId, pushErr.Error())
		log.Error(e.Error())
		return e
	}
	
	return nil
}

func (sdp *scyllaDbPusher) Close() error {
	return sdp.pool.Close()
}