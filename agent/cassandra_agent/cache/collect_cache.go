package cache

import (
	"agent_common/pkg/collector/host"
)

type _CollectCache struct {
	Agent struct {
		Cpu struct {
			Data host.HostCpuStat
			Time int64
		}
	}
}

var (
	CollectCache = _CollectCache{}
)