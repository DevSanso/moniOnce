package pusher

import (
	"cassandra_agent/types"
	"context"
)

func pushHostCpu(objectId int, ctx context.Context, conn IConnPusher, root *types.PushData, ttl int64) error {
	data := &root.Host.CpuPercent
	return conn.Exec(ctx, _InsertHostCpuQuery, objectId, root.NowTimeUnixEpoch, data.User, data.System, data.Wait, data.Idle, ttl)
}

func pushHostMemory(objectId int, ctx context.Context, conn IConnPusher, root *types.PushData, ttl int64) error {
	data := &root.Host.Memory
	return conn.Exec(ctx, _InsertHostCpuQuery, objectId, root.NowTimeUnixEpoch, data.Total, data.Use, data.Free, ttl)
}

func pushClients(objectId int, ctx context.Context, conn IConnPusher, root *types.PushData, ttl int64) error {
	data := root.Cql.Clients
	var err error = nil

	for _, c := range data {
		err = conn.Exec(ctx, _InsertCQLClientsQuery, objectId, root.NowTimeUnixEpoch, c.Address, c.DriverName, c.ConnectionStage, c.Hostname, c.RequestCnt, c.Keyspace, c.Username, ttl)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func pushRunningQuery(objectId int, ctx context.Context, conn IConnPusher, root *types.PushData, ttl int64) error {
	data := root.Cql.RunningQuery
	var err error = nil

	for _,q := range data {
		err = conn.Exec(ctx, _InsertCQLRunningQueryQuery, objectId, root.NowTimeUnixEpoch, q.ThreadId, q.QueueMicroSec, q.RunningMicroSec, q.Text)
		if err != nil {
			return err
		}
	}

	return nil
}