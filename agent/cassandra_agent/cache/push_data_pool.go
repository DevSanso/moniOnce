package cache

import (
	"cassandra_agent/types/dataframe"
	"sync"
)

var (
	NodetoolTpStatMemoryPool sync.Pool = sync.Pool {
		New : func() any {
			l := make([]dataframe.LatencyMetrics, 20)
			p := make([]dataframe.PoolMetrics, 20)

			return struct {
				l []dataframe.LatencyMetrics;
				p []dataframe.PoolMetrics
			}{l, p }
		},
	}

	NodetoolInfoMemoryPool sync.Pool = sync.Pool {
		New : func() any {
			return &dataframe.InfoMetrics{}
		},
	}
)