package main

import (
	"agent_common/pkg/applnew"
	appltype "agent_common/pkg/applnew/types"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cassandra_agent/cassandra"
	"cassandra_agent/collect"
	agenttype "cassandra_agent/types"
)


func main() {
	args := NewArgs()

	application := applnew.NewApplication[agenttype.PushData, *cassandra.CassandraConn, agenttype.FlagData]()
	err := application.Init(appltype.InitData[agenttype.PushData, *cassandra.CassandraConn, agenttype.FlagData, *agenttype.FlagData]{
		SettingPath: args.configPath, CollectM: collect.CollectMapping, CronM: nil, DataPusher: nil, GetConnPoolFn: cassandra.NewCassandraPool,
	})

	if err != nil {
		log.Println("init failed : ", err.Error())
		return
	}

	applCtx, cancelCtxFn := context.WithCancel(context.Background())

	if runErr := application.Run(applCtx); runErr != nil {
		log.Println("Main - Server Run Error : ", runErr.Error())
		return
	}

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT)

	isStop := false
	for isStop {
		select {
		case <-osSignal:
			cancelCtxFn()
			time.Sleep(5 * time.Second)
			log.Println("Main - stop server")
			isStop = true
		default:
		}
		time.Sleep(3 * time.Second)
	}
}