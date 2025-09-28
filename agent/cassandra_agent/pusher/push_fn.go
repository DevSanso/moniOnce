package pusher

import (
	"cassandra_agent/types/dataframe"
	"context"
)

func pushHostCpu(ctx context.Context, conn IConnPusher, data *dataframe.HostCpuPercent) error {
	return conn.Exec(ctx, _InsertHostCpuQuery, data.User, data.System, data.Wait, data.Idle)
}

func pushHostMemory(ctx context.Context, conn IConnPusher, data *dataframe.HostMemory) error {
	return conn.Exec(ctx, _InsertHostCpuQuery, data.Total, data.Use, data.Free)
}

func pushClients(ctx context.Context, conn IConnPusher, data []dataframe.SystemViewClients) error {
	var err error = nil

	for _, c := range data {
		err = conn.Exec(ctx, _InsertCQLClientsQuery, c.Address, c.DriverName, c.ConnectionStage, c.Hostname, c.RequestCnt, c.Keyspace, c.Username)
		if err != nil {
			return err
		}
	}
	
	return nil
}