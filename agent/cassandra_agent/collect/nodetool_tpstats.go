package collect

import (
	"agent_common/pkg/applnew/logger"
	"bytes"
	"cassandra_agent/cache"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func collectNodetoolTpStats(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	stdout.Grow(2048)

	cmd := exec.Command("nodetool", "tpstats")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if cmdErr := stderr.String(); cmdErr != "" {
		log.Error(cmdErr)
		return nil, fmt.Errorf("%s", cmdErr)
	}

	var poolData string = ""
	var LatencyData string = ""

	parts := strings.Split(stdout.String(), "\n\n")
	if len(parts) >= 2 {
		poolData = parts[0]
		LatencyData = parts[1]
	}

	cache := cache.NodetoolTpStatMemoryPool.Get().(struct {
		l []dataframe.LatencyMetrics
		p []dataframe.PoolMetrics
	})

	poolMetric := cache.p
	latencyMetric := cache.l

	poolErr := dataframe.ParsePoolMetrics(poolData, poolMetric)
	if poolErr != nil {
		log.Error(poolErr.Error())
		return nil, poolErr
	}
	latencyErr := dataframe.ParseLatencyMetrics(LatencyData, latencyMetric)
	if latencyErr != nil {
		log.Error(latencyErr.Error())
		return nil, latencyErr
	}

	pushData := new(types.PushData)
	pushData.NowTimeUnixEpoch = time.Now().Unix()
	pushData.ConnTypeId = int(constants.ConnTypeNodeTool)
	pushData.DataId = int(constants.NodeToolTpStatsData)
	pushData.Nodetool.TpStats.Pool = poolMetric
	pushData.Nodetool.TpStats.Latency = latencyMetric

	return pushData, nil
}
