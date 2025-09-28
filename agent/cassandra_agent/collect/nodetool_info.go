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
	"time"
)

func collectNodeToolInfo(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	stdout.Grow(2048)

	cmd := exec.Command("nodetool", "info")
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
	data := cache.NodetoolInfoMemoryPool.Get().(*dataframe.InfoMetrics)

	parserErr := dataframe.ParseInfoMetrics(stdout.String(), data)
	if parserErr != nil {
		log.Error(parserErr.Error())
		return nil, fmt.Errorf("%s", parserErr)
	}

	pushData := new(types.PushData)
	pushData.NowTimeUnixEpoch = time.Now().Unix()
	pushData.ConnTypeId = int(constants.ConnTypeNodeTool)
	pushData.DataId = int(constants.NodeToolInfoData)
	pushData.Nodetool.Info = data

	return pushData, nil
}
