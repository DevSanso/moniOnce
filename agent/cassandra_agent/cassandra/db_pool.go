package cassandra

import (
	apptype "agent_common/pkg/applnew/types"
	"context"
	"fmt"

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

// GetDbConn implements types.CollectConnPool.
func (cp *CassandraPool) GetDbConn(ctx context.Context) (*CassandraConn, error) {
	panic("unimplemented")
}

type CassandraConn struct {
	cqlConfig *gocql.ClusterConfig
	cqlSession *gocql.Session
}

func (cc *CassandraConn) Close() error {
	if cc.cqlSession != nil {
		cc.cqlSession.Close()
		cc.cqlSession = nil
	}

	return nil
}

func (cc *CassandraConn) ConnectCQL() error {
	if cc.cqlSession != nil {
		return fmt.Errorf("already new cqlSession")
	}

	session, sessErr := cc.cqlConfig.CreateSession()
	if sessErr != nil {
		return sessErr
	}
	cc.cqlSession = session

	return nil
}

func (cc *CassandraConn) Query(ctx context.Context, query string, genFn func(scanFn func(...any) error) (any, error), args ...any) ([]any, error) {
	q := cc.cqlSession.Query(query, args...)
	iter := q.IterContext(ctx)
	scanner := iter.Scanner()

	ret := make([]any, 5)
	
	for scanner.Next() {
		data, err := genFn(scanner.Scan)
		if err != nil {
			iter.Close()
			return nil, err
		}
		ret = append(ret, data)
	}
	iter.Close()

	return ret, nil
}

func (cp *CassandraPool) Close() error {
	return nil
}

func NewCassandraPool(IP string, Port int, User string, Password string, Dbname string, args ...any) (apptype.CollectConnPool[*CassandraConn], error) {
	cconf := gocql.NewCluster(fmt.Sprintf("%s:%d", IP, Port))
	password := gocql.PasswordAuthenticator{
		Username: User,
		Password: Password,
	}
	cconf.Keyspace = Dbname
	cconf.Authenticator = password

	cp := &CassandraPool{
		connectInfo: cconf,
	}

	return cp, nil
}
