package cassandra

import (
	apptype "agent_common/pkg/applnew/types"
	"context"
	"fmt"
	"net/url"

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

func (cc *CassandraConn) connectCQL() error {
	if cc.cqlSession != nil {
		return nil
	}

	session, sessErr := cc.cqlConfig.CreateSession()
	if sessErr != nil {
		return sessErr
	}
	cc.cqlSession = session

	return nil
}

func (cc *CassandraConn) Exec(ctx context.Context, query string, args ...any) error {
	if err := cc.connectCQL(); err != nil {
		return err
	}

	q := cc.cqlSession.Query(query, args...)
	if err := q.ExecContext(ctx); err != nil {
		return err
	}
	
	return nil
}

func CassandraConnRunQuery[T any](cc *CassandraConn, ctx context.Context, query string, cap int, genFn func(p *T, scanFn func(...any) error) error, args ...any) ([]T, error){
	if err := cc.connectCQL(); err != nil {
		return nil, err
	}
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

func NewCassandraPool(dbUrl string, args ...any) (apptype.CollectConnPool[*CassandraConn], error) {
	url, err := url.Parse(dbUrl)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "cassandra_agent" {
		return nil, fmt.Errorf("not support schema : %s", url.Scheme)
	}

	ip := url.Host
	port := url.Port()

	user := url.User.Username()
	passwd,_ := url.User.Password()
	dbname := url.Path

	cconf := gocql.NewCluster(fmt.Sprintf("%s:%s", ip, port))
	password := gocql.PasswordAuthenticator{
		Username: user,
		Password: passwd,
	}
	cconf.Keyspace = dbname
	cconf.Authenticator = password

	cp := &CassandraPool{
		connectInfo: cconf,
	}

	return cp, nil
}
