package collect

import (
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/collector/host"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
	"time"
)

func collectAgentHostMem(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	stat, statErr := host.NewHostMemoryCollector().HostMemory()

	if statErr != nil {
		log.Error(statErr.Error())
		return nil, statErr
	}

	pushData := new(types.PushData)
	pushData.NowTimeUnixEpoch = time.Now().Unix()
	pushData.ConnTypeId = int(constants.ConnTypeAgent)
	pushData.DataId = int(constants.AgentHostMemory)
	pushData.Host.Memory = dataframe.HostMemory(stat)

	return pushData, nil
}
