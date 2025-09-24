package collect

import (
	"agent_common/pkg/applnew/logger"
	"bytes"
	"cassandra_agent/cassandra"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
	"os/exec"
	"strings"
)

func CollectNodeToolTpStats(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error){
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	
	stdout.Grow(1024)

	cmd := exec.Command("nodetool", "tpstats")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var poolData string = ""
	var LatencyData string = ""

	parts := strings.Split(stdout.String(), "\n\n")
	if len(parts) >= 2 {
		poolData = parts[0]
		LatencyData =  parts[1]
	}

	poolMetric, poolErr := dataframe.ParsePoolMetrics(poolData)
	if poolErr != nil {
		return nil, poolErr
	}
	latencyMetric, latencyErr := dataframe.ParseLatencyMetrics(LatencyData)
	if latencyErr != nil {
		return nil, latencyErr
	}

	ret := new(types.PushData)
	ret.ConnType = "NODETOOL"
	ret.DataName = "TPSTATS"
	ret.Nodetool.TpStats.Pool = poolMetric
	ret.Nodetool.TpStats.Latency = latencyMetric

	return ret, nil
}