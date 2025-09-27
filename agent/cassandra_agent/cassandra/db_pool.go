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

func CassandraConnRunQuery[T any](cc *CassandraConn, ctx context.Context, query string, cap int, 
	genFn func(p *T, scanFn func(...any) error) error, args ...any) ([]T, error){
	q := cc.cqlSession.Query(query, args...)
	iter := q.IterContext(ctx)
	scanner := iter.Scanner()

	rows := make([]T, 0, cap)
	for scanner.Next() {
		var row T
		err := genFn(&row, scanner.Scan)
		if err != nil {
			iter.Close()
			return nil, err
		}

		rows = append(rows, row)
	}
	iter.Close()

	return rows, nil
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
