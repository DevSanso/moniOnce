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

func CollectAgentHostCpu(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error){
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
		pushData.ConnTypeId = int(constants.ConnTypeAgent)
		pushData.DataId = int(constants.AgentHostCpu)
		pushData.Agent.CpuPercent = dataframe.AgentHostCpuPercent{}
		return pushData, nil
	}

	totalDelta := float64(totalTime - cache.CollectCache.Agent.Cpu.Time)
	sysDelta   := stat.Sys - cache.CollectCache.Agent.Cpu.Data.Sys
	userDelta  := stat.User - cache.CollectCache.Agent.Cpu.Data.User
	waitDelta := stat.Wait - cache.CollectCache.Agent.Cpu.Data.Wait
	idleDelta := stat.Idle - cache.CollectCache.Agent.Cpu.Data.Idle

	cache.CollectCache.Agent.Cpu.Data = stat
	cache.CollectCache.Agent.Cpu.Time = totalTime

	pushData := new(types.PushData)
	pushData.ConnTypeId = int(constants.ConnTypeAgent)
	pushData.DataId = int(constants.AgentHostCpu)
	pushData.Agent.CpuPercent = dataframe.AgentHostCpuPercent{
		System : sysDelta / totalDelta,
		User   : userDelta / totalDelta,
		Wait   : waitDelta / totalDelta,
		Idle   : idleDelta / totalDelta,
	}

	return pushData, nil	
}