package collect

import (
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/collector/host"
	"cassandra_agent/cache"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
	"time"
)

func collectAgentHostCpu(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	stat, statErr := host.NewHostCpuCollector().HostCpu()
	totalTime := time.Now().UnixMilli()

	if statErr != nil {
		log.Error(statErr.Error())
		return nil, statErr
	}

	if cache.CollectCache.Agent.Cpu.Time == 0 {
		log.Debug("CollectAgentHostCpu - first init cpu data")

		cache.CollectCache.Agent.Cpu.Data = stat
		cache.CollectCache.Agent.Cpu.Time = totalTime

		pushData := new(types.PushData)
		pushData.NowTimeUnixEpoch = int64(totalTime / 1000)
		pushData.ConnTypeId = int(constants.ConnTypeAgent)
		pushData.DataId = int(constants.AgentHostCpu)
		pushData.Host.CpuPercent = dataframe.HostCpuPercent{}
		return pushData, nil
	}

	totalDelta := float64(totalTime - cache.CollectCache.Agent.Cpu.Time)
	sysDelta := stat.Sys - cache.CollectCache.Agent.Cpu.Data.Sys
	userDelta := stat.User - cache.CollectCache.Agent.Cpu.Data.User
	waitDelta := stat.Wait - cache.CollectCache.Agent.Cpu.Data.Wait
	idleDelta := stat.Idle - cache.CollectCache.Agent.Cpu.Data.Idle

	cache.CollectCache.Agent.Cpu.Data = stat
	cache.CollectCache.Agent.Cpu.Time = totalTime

	pushData := new(types.PushData)
	pushData.NowTimeUnixEpoch = int64(totalTime / 1000)
	pushData.ConnTypeId = int(constants.ConnTypeAgent)
	pushData.DataId = int(constants.AgentHostCpu)
	pushData.Host.CpuPercent = dataframe.HostCpuPercent{
		System: sysDelta / totalDelta,
		User:   userDelta / totalDelta,
		Wait:   waitDelta / totalDelta,
		Idle:   idleDelta / totalDelta,
	}

	return pushData, nil
}
