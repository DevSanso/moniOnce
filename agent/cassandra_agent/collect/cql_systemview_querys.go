package collect

import (
	"agent_common/pkg/applnew/logger"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
	"time"
)

func collectCQLSystemViewRunningQuery(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	rows, rowsErr := cassandra.CassandraConnRunQuery(ctl, ctx, _CqlSystemViewQueriesQuery, 5, func(p *dataframe.SystemViewQueries, scanFn func(...any) error) error {
		row := p
		if err := scanFn(&row.ThreadId, &row.QueueMicroSec, &row.RunningMicroSec, &row.Text); err != nil {
			return err
		}

		return nil
	})

	if rowsErr != nil {
		log.Error(rowsErr.Error())
		return nil, rowsErr
	}

	pushData := new(types.PushData)
	pushData.NowTimeUnixEpoch = time.Now().Unix()
	pushData.ConnTypeId = int(constants.ConnTypeCQLTool)
	pushData.DataId = int(constants.CQLRunningQuerys)
	pushData.Cql.RunningQuery = rows

	return pushData, nil

}
