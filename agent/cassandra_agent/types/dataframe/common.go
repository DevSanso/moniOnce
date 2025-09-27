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