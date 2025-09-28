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

func collectCQLSystemViewClients(ctx context.Context, ctl *cassandra.CassandraConn, log logger.LevelLogger) (*types.PushData, error) {
	rows, rowsErr := cassandra.CassandraConnRunQuery(ctl, ctx, _CqlSystemViewClientsQuery, 5, func(p *dataframe.SystemViewClients, scanFn func(...any) error) error {
		row := p
		if err := scanFn(&row.Address, &row.DriverName, &row.ConnectionStage, &row.Hostname, &row.RequestCnt, &row.Keyspace, &row.Username); err != nil {
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
	pushData.DataId = int(constants.CQLClients)
	pushData.Cql.Clients = rows

	return pushData, nil

}
