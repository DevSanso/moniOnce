package host

type HostMemoryStat struct {
	Use uint64
	Total  uint64
	Free   uint64
}

type HostCpuStat struct {
	Sys float64
	User float64
	Wait float64
	Idle float64
}