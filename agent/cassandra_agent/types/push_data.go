package types

import (
	"cassandra_agent/types/dataframe"
)

type PushData struct {
	ConnTypeId int
	DataId     int

	NowTimeUnixEpoch int64    

	Nodetool struct {
		TpStats struct {
			Pool    []dataframe.PoolMetrics
			Latency []dataframe.LatencyMetrics
		}

		Info *dataframe.InfoMetrics
	}

	Cql struct {
		TracesSession []dataframe.TracesSession
		RunningQuery  []dataframe.SystemViewQueries
		Clients       []dataframe.SystemViewClients
		SystemLogs    []dataframe.SystemViewSystemLog
	}

	Host struct {
		CpuPercent dataframe.HostCpuPercent
		Memory     dataframe.HostMemory
	}
}
