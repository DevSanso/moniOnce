package host

func NewHostMemoryCollector() HostMemoryCollector {
	var o gopsHostMemoryCollector
	return o
}

func NewHostCpuCollector() HostCpuCollector {
	var o gopsHostCpuCollector
	return o
}