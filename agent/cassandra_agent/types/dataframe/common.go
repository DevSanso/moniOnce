package dataframe

type PoolMetrics struct {
	PoolName       string `agent_common_parser:"0,string"`
	Active         int `agent_common_parser:"1,int"`
	Pending        int `agent_common_parser:"2,int"`
	Completed      int `agent_common_parser:"3,int"`
	Blocked        int `agent_common_parser:"4,int"`
	AllTimeBlocked int `agent_common_parser:"5,int"`
}

type CacheMetric struct {
	Entries            int
	SizeByte           int64
	CapacityByte       int64
	Hits               int
	Requests           int
	HitRate            float64
	SavePeriodInSecond int
	OverflowSizeByte   int64
}

type LatencyMetrics struct {
	MessageType      string  `agent_common_parser:"0,string"`
	Dropped          int     `agent_common_parser:"1,int"`
	Latency50Percent float64 `agent_common_parser:"2,float64"`
	Latency95Percent float64 `agent_common_parser:"3,float64"`
	Latency99Percent float64 `agent_common_parser:"4,float64"`
	LatencyMax       float64 `agent_common_parser:"5,float64"`
}
