package types

import "cassandra_agent/types/dataframe"

type PushData struct {
	ConnTypeId int
	DataId int
	
	Nodetool struct {
		TpStats struct {
			Pool []dataframe.PoolMetrics
			Latency []dataframe.LatencyMetrics
		}

		Info *dataframe.InfoMetrics
	}

	Cql struct {
		TracesSession []dataframe.TracesSession
		RunningQuery  []dataframe.SystemViewQueries
	}

	Agent struct {
		CpuPercent    dataframe.AgentHostCpuPercent
		Memory        dataframe.AgentHostMemory
	}
}