package collect

import (
	"agent_common/pkg/applnew/constants/context/key"
	"agent_common/pkg/applnew/logger"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
	"time"
)

func collectCQLSystemViewSystemLogs(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	interval := ctx.Value(key.ContextIntervalKey).(int)
	
	rows, rowsErr := cassandra.CassandraConnRunQuery(ctl, ctx, _CqlSystemViewSystemLogsQuery, 5, func(p *dataframe.SystemViewSystemLog, scanFn func(...any) error) error {
		row := p
		if err := scanFn(&row.Timestamp, &row.Level, &row.Logger, &row.Message); err != nil {
			return err
		}

		return nil
	}, interval)

	if rowsErr != nil {
		log.Error(rowsErr.Error())
		return nil, rowsErr
	}

	pushData := new(types.PushData)
	pushData.NowTimeUnixEpoch = time.Now().Unix()
	pushData.ConnTypeId = int(constants.ConnTypeCQLTool)
	pushData.DataId = int(constants.CQLSystemLog)
	pushData.Cql.SystemLogs = rows

	return pushData, nil

}