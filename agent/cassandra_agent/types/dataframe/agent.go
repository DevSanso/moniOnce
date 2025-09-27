package dataframe

import "agent_common/pkg/collector/host"

type AgentHostCpuPercent struct {
	System float64
	User   float64
	Wait   float64
	Idle   float64
}
type AgentHostMemory host.HostMemoryStat