package host

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type gopsHostMemoryCollector int
type gopsHostCpuCollector int

func (gmc gopsHostMemoryCollector) HostMemory() (HostMemoryStat, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return HostMemoryStat{}, err
	}

	return HostMemoryStat{
		Free: mem.Free,
		Use: mem.Used,
		Total: mem.Total,
	}, nil
} 

func (gmc gopsHostCpuCollector) HostCpu() (HostCpuStat, error){
	tempCpuTime, err := cpu.Times(false)
	if err != nil {
		return HostCpuStat{}, err
	}
	cpuTime := &tempCpuTime[0]
	
	return HostCpuStat{
		Sys : cpuTime.System,
		User : cpuTime.User,
		Wait : cpuTime.Iowait,
		Idle : cpuTime.Idle,
	}, nil	
}


