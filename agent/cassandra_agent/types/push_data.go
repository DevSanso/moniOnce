package types

import "cassandra_agent/types/dataframe"

type PushData struct {
	ConnType string
	DataName string
	
	Nodetool struct {
		TpStats struct {
			Pool []dataframe.PoolMetrics
			Latency []dataframe.LatencyMetrics
		}
	}
}