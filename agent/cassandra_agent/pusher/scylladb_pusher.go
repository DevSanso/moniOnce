package pusher

import (
	"agent_common/pkg/applnew/logger"
	apptype "agent_common/pkg/applnew/types"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
	"context"
)

type scyllaDbPusher struct {
	pool apptype.CollectConnPool[*cassandra.CassandraConn]
}

func NewScyllaDbPusher(IP string, Port int, User string, Password string, Dbname string, args ...any) (apptype.DataPusher[types.PushData], error) {
	p, err := cassandra.NewCassandraPool(IP, Port, User, Password, Dbname, args...)
	if err != nil {
		return nil, err
	}
	return &scyllaDbPusher{pool : p}, nil
}


func (sdp *scyllaDbPusher) Push(data *types.PushData, ctx context.Context, log logger.LevelLogger) error {

	return nil
}

func (sdp *scyllaDbPusher) Close() error {

	return nil
}