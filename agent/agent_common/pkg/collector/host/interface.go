package host

type HostMemoryCollector interface {
	HostMemory() (HostMemoryStat, error)
}

type HostCpuCollector interface {
	HostCpu() (HostCpuStat, error)
}
