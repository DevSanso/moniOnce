package cassandra

import (
	"agent_common/pkg/logger"
	"agent_common/pkg/util/types"
	"cassandra_agent/config"
	"context"
	"io"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)


type CassandraConnType string

const (
	CassandraCmdConnType CassandraConnType = "cmd"
	CassandraCQLConnType CassandraConnType = "cql"
)

type CassandraPool struct {
	connectInfo *gocql.ClusterConfig


}

func newCassandraPool(conf *config.DBConfig) (*CassandraPool, io.Closer, error) {
	cconf := gocql.NewCluster(conf.Cassnadra.Addr)
	password := gocql.PasswordAuthenticator {
		Username: conf.Cassnadra.User,
		Password: conf.Cassnadra.Password,
	}

	cconf.Authenticator = password

	cp := &CassandraPool {
		connectInfo: cconf,
	}

	return cp, cp, nil
}

func NewCassandraPool(conf *config.DBConfig) (types.TargetDbPool[CassandraCollectConnCtl], io.Closer, error) {
	return newCassandraPool(conf)
}

func NewCassandraDbLogger(conf *config.DBConfig) (logger.DbLogger, io.Closer, error) {
	return newCassandraPool(conf)
}

func (cp *CassandraPool) Exec(query string, args [][]any) error {
	return nil	
}

func (cp *CassandraPool)GetDbConn(ctx context.Context) (CassandraCollectConnCtl, error) {
	return CassandraCollectConnCtl{}, nil	
}

func (cp *CassandraPool)Close() error {
	return nil
}

type CassandraCollectConnCtl struct {
	connectInfo *gocql.ClusterConfig
}

func (ccct *CassandraCollectConnCtl) GetCQLConnection() {

}

func (ccct *CassandraCollectConnCtl) GetNodeToolCmdConnection() {
	
}