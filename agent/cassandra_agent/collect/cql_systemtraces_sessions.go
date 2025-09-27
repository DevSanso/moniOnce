package collect

import (
	"agent_common/pkg/applnew/logger"
	"cassandra_agent/cassandra"
	"cassandra_agent/constants"
	"cassandra_agent/types"
	"cassandra_agent/types/dataframe"
	"context"
)

func CollectCQLSystemTracesSessions(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error){
	if err := ctl.ConnectCQL(); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	rows, rowsErr := cassandra.CassandraConnRunQuery(ctl, ctx, _CqlSystemTracesQuery, 5, func(p *dataframe.TracesSession, scanFn func(...any) error) error {
		row := p
		if err := scanFn(&row.SessionID, &row.Client, &row.Command, &row.Coordinator, &row.CoordiantorPort, &row.Duration, &row.Parameters, &row.Request, &row.Started_at); err != nil {
			return err
		}
	
		return nil
	})

	if rowsErr != nil {
		log.Error(rowsErr.Error())
		return nil, rowsErr
	}

	pushData := new(types.PushData)
	pushData.ConnTypeId = int(constants.ConnTypeCQLTool)
	pushData.DataId = int(constants.CQLTracesSessions)
	pushData.Cql.TracesSession = rows

	return pushData, nil

}