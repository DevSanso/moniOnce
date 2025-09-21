package main

import (
	"agent_common/pkg/appl"
	"agent_common/pkg/util/types"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cassandra_agent/cassandra"
	"cassandra_agent/collect"
	"cassandra_agent/config"
)

func main() {
	args := NewArgs()
	
	extend := appl.AgentApplicationExtendInitConfig[cassandra.CassandraCollectConnCtl, *config.DBConfig] {
		GenTargetDbPoolrFn: cassandra.NewCassandraPool,
		GenDbLoggerFn: cassandra.NewCassandraDbLogger,
		DataLoggers: nil,
		Intervals: []types.IntervalRegister[cassandra.CassandraCollectConnCtl]{
			{Name : "cql.system.local", Fn : collect.CollectSystemLocalhandle},
		},
	}

	server, err := appl.InitAgentApplication(args.configPath, extend, nil)
	if err != nil {
		log.Println("Main - Agent Application Init Failed :", err.Error())
		return
	}

	applCtx, cancelCtxFn := context.WithCancel(context.Background())

	if runErr := server.Run(applCtx); runErr != nil {
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