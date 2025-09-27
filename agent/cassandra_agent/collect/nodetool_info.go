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
)

func CollectNodeToolInfo(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error){
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

	ret := new(types.PushData)
	ret.ConnTypeId = int(constants.ConnTypeNodeTool)
	ret.DataId = int(constants.NodeToolInfoData)
	ret.Nodetool.Info = data

	return ret, nil
}