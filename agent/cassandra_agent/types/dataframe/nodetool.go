package dataframe

import (
	"agent_common/pkg/util/funcs"
	"agent_common/pkg/util/parser"
	"fmt"
	"strconv"
	"strings"
)

type PoolMetrics struct {
	PoolName       string `agent_common_parser:"0,string"`
	Active         int `agent_common_parser:"1,int"`
	Pending        int `agent_common_parser:"2,int"`
	Completed      int `agent_common_parser:"3,int"`
	Blocked        int `agent_common_parser:"4,int"`
	AllTimeBlocked int `agent_common_parser:"5,int"`
}

type LatencyMetrics struct {
	MessageType      string `agent_common_parser:"0,string"`
	Dropped          int `agent_common_parser:"1,int"`
	Latency50Percent float64 `agent_common_parser:"2,float64"`
	Latency95Percent float64 `agent_common_parser:"3,float64"`
	Latency99Percent float64 `agent_common_parser:"4,float64"`
	LatencyMax       float64 `agent_common_parser:"5,float64"`
}

type InfoMetricsCache struct {
	Entries int
	SizeByte  int64
	CapacityByte int64
	Hits         int
	Requests     int
	HitRate      float64
	SavePeriodInSecond int
	OverflowSizeByte   int64
}
type InfoMetrics struct {
	ID string
	IsGossipActive bool
	IsNativeTransportActive bool
	LoadKb float64
	UncompressedLoadKb float64
	GenerationNo int64
	Uptime       int64
	HeapMemUsageMB float64
	HeapMemTotalMB float64
	OffHeapMemMB   float64
	DataCenter     string
	Rack           string
	Exceptions     int
	Key            InfoMetricsCache
	Row            InfoMetricsCache
	Counter        InfoMetricsCache
	Network        InfoMetricsCache
	RepairedPercent int
	BootstrapState string
}

var (
	poolMetricParser = parser.CreateLinePaser[PoolMetrics](" ")
	latencyMetricParser = parser.CreateLinePaser[LatencyMetrics](" ")
)

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

func ParsePoolMetrics(data string, poolMetricsArray []PoolMetrics) error {
	if maxLen := countPoolMetrics(data); maxLen > len(poolMetricsArray) {
		poolMetricsArray = nil
		poolMetricsArray = make([]PoolMetrics, maxLen)
	} 

	lines := strings.Split(data, "\n")
	poolIndex := 0
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if err := poolMetricParser.Load(line, &poolMetricsArray[poolIndex]); err != nil {
			return fmt.Errorf("error initializing PoolMetrics: %v", err)
		}
		poolIndex++
	}
	return nil
}

func ParseLatencyMetrics(data string, latencyMetricsArray []LatencyMetrics) error {
	if maxLen := countLatencyMetrics(data); maxLen > len(latencyMetricsArray) {
		latencyMetricsArray = nil
		latencyMetricsArray = make([]LatencyMetrics, maxLen)
	} 

	lines := strings.Split(data, "\n")
	latencyIndex := 0
	for _, line := range lines[1:] {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if err := latencyMetricParser.Load(line, &latencyMetricsArray[latencyIndex]); err != nil {
			return fmt.Errorf("error initializing LatencyMetrics: %v", err)
		}
		latencyIndex++
	}
	return nil
}

func parseCacheLine(value string) InfoMetricsCache {
	fields := strings.Split(value, ",")
	cache := InfoMetricsCache{}
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if strings.HasPrefix(f, "entries ") {
			cache.Entries, _ = strconv.Atoi(strings.TrimPrefix(f, "entries "))
		} else if strings.HasPrefix(f, "size ") {
			cache.SizeByte = funcs.ParseBytesFromStr(strings.TrimPrefix(f, "size "))
		} else if strings.HasPrefix(f, "overflow size: ") {
			cache.OverflowSizeByte = funcs.ParseBytesFromStr(strings.TrimPrefix(f, "overflow size: "))
		} else if strings.HasPrefix(f, "capacity ") {
			cache.CapacityByte = funcs.ParseBytesFromStr(strings.TrimPrefix(f, "capacity "))
		} else if strings.HasSuffix(f, "hits") {
			cache.Hits, _ = strconv.Atoi(strings.TrimSuffix(f, " hits"))
		} else if strings.HasSuffix(f, "requests") {
			cache.Requests, _ = strconv.Atoi(strings.TrimSuffix(f, " requests"))
		} else if strings.HasSuffix(f, "recent hit rate") {
			v := strings.TrimSuffix(f, " recent hit rate")
			cache.HitRate, _ = strconv.ParseFloat(v, 64)
		} else if strings.HasSuffix(f, "save period in seconds") {
			v := strings.TrimSuffix(f, " save period in seconds")
			cache.SavePeriodInSecond, _ = strconv.Atoi(v)
		}
	}
	return cache
}

func ParseInfoMetrics(data string, info *InfoMetrics) (error) {
	ret := info
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ID":
			ret.ID = value
		case "Gossip active":
			ret.IsGossipActive = funcs.ParseBool(value)
		case "Native Transport active":
			ret.IsNativeTransportActive = funcs.ParseBool(value)
		case "Load":
			ret.LoadKb = float64(funcs.ParseBytesFromStr(value)) / 1024
		case "Uncompressed load":
			ret.UncompressedLoadKb = float64(funcs.ParseBytesFromStr(value)) / 1024
		case "Generation No":
			ret.GenerationNo, _ = strconv.ParseInt(value, 10, 64)
		case "Uptime (seconds)":
			ret.Uptime, _ = strconv.ParseInt(value, 10, 64)
		case "Heap Memory (MB)":
			memParts := strings.Split(value, "/")
			if len(memParts) == 2 {
				ret.HeapMemUsageMB, _ = strconv.ParseFloat(strings.TrimSpace(memParts[0]), 64)
				ret.HeapMemTotalMB, _ = strconv.ParseFloat(strings.TrimSpace(memParts[1]), 64)
			}
		case "Off Heap Memory (MB)":
			ret.OffHeapMemMB, _ = strconv.ParseFloat(value, 64)
		case "Data Center":
			ret.DataCenter = value
		case "Rack":
			ret.Rack = value
		case "Exceptions":
			ret.Exceptions, _ = strconv.Atoi(value)
		case "Key Cache":
			ret.Key = parseCacheLine(value)
		case "Row Cache":
			ret.Row = parseCacheLine(value)
		case "Counter Cache":
			ret.Counter = parseCacheLine(value)
		case "Network Cache":
			ret.Network = parseCacheLine(value)
		case "Percent Repaired":
			val := strings.TrimSuffix(value, "%")
			f, _ := strconv.ParseFloat(val, 64)
			ret.RepairedPercent = int(f)
		case "Bootstrap state":
			ret.BootstrapState = value
		}
	}

	return nil
}