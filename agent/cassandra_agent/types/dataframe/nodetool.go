package dataframe

import (
	"fmt"
	"strings"
	"strconv"
)

type PoolMetrics struct {
	PoolName       string
	Active         int
	Pending        int
	Completed      int
	Blocked        int
	AllTimeBlocked int
}

type LatencyMetrics struct {
	MessageType      string
	Dropped          int
	Latency50Percent float64
	Latency95Percent float64
	Latency99Percent float64
	LatencyMax       float64
}

func initLatencyMetrics(ptr *LatencyMetrics, data string) error {
	fields := strings.Fields(data)
	if len(fields) != 6 {
		return fmt.Errorf("invalid LatencyMetrics data: expected 6 fields, got %d", len(fields))
	}

	ptr.MessageType = fields[0]
	var err error
	ptr.Dropped, err = strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid Dropped value: %v", err)
	}
	ptr.Latency50Percent, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return fmt.Errorf("invalid Latency50Percent value: %v", err)
	}
	ptr.Latency95Percent, err = strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return fmt.Errorf("invalid Latency95Percent value: %v", err)
	}
	ptr.Latency99Percent, err = strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return fmt.Errorf("invalid Latency99Percent value: %v", err)
	}
	ptr.LatencyMax, err = strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return fmt.Errorf("invalid LatencyMax value: %v", err)
	}

	return nil
}

func initPoolMetrics(ptr *PoolMetrics, data string) error {
	fields := strings.Fields(data)
	if len(fields) != 6 {
		return fmt.Errorf("invalid PoolMetrics data: expected 6 fields, got %d", len(fields))
	}

	ptr.PoolName = fields[0]
	var err error
	ptr.Active, err = strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid Active value: %v", err)
	}
	ptr.Pending, err = strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("invalid Pending value: %v", err)
	}
	ptr.Completed, err = strconv.Atoi(fields[3])
	if err != nil {
		return fmt.Errorf("invalid Completed value: %v", err)
	}
	ptr.Blocked, err = strconv.Atoi(fields[4])
	if err != nil {
		return fmt.Errorf("invalid Blocked value: %v", err)
	}
	ptr.AllTimeBlocked, err = strconv.Atoi(fields[5])
	if err != nil {
		return fmt.Errorf("invalid AllTimeBlocked value: %v", err)
	}

	return nil
}

func countPoolMetrics(data string) int {
	lines := strings.Split(data, "\n")
	count := 0
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		count++
	}
	return count
}

func countLatencyMetrics(data string) int {
	lines := strings.Split(data, "\n")
	count := 0
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		count++
	}
	return count
}

func ParsePoolMetrics(data string) ([]PoolMetrics, error) {
	poolMetricsArray := make([]PoolMetrics, countPoolMetrics(data))

	lines := strings.Split(data, "\n")
	poolIndex := 0
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if err := initPoolMetrics(&poolMetricsArray[poolIndex], line); err != nil {
			return nil, fmt.Errorf("error initializing PoolMetrics: %v", err)
		}
		poolIndex++
	}
	return poolMetricsArray, nil
}

func ParseLatencyMetrics(data string) ([]LatencyMetrics,error) {
	latencyMetricsArray := make([]LatencyMetrics, countLatencyMetrics(data))

	lines := strings.Split(data, "\n")
	latencyIndex := 0
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if err := initLatencyMetrics(&latencyMetricsArray[latencyIndex], line); err != nil {
			return nil, fmt.Errorf("error initializing LatencyMetrics: %v", err)
		}
		latencyIndex++
	}
	return latencyMetricsArray, nil
}