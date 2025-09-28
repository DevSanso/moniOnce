package dataframe

import "agent_common/pkg/collector/host"

type HostCpuPercent struct {
	System float64
	User   float64
	Wait   float64
	Idle   float64
}
type HostMemory host.HostMemoryStat
