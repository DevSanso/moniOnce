package config

import (
	"agent_common/internal/constant"
	"agent_common/pkg/logger"
	"agent_common/pkg/util/types"
	"context"
	"database/sql"
)

type sqlConfigure[T types.StrSetter] struct {
	c * sql.Conn
	objectId int
	marker *T
	levelLogger logger.LevelLogger
}

func NewSQLConfigure[T types.StrSetter](conn *sql.Conn, objectId int, levelLogger logger.LevelLogger) Configure[T] {
	return &sqlConfigure[T] {
		c: conn,
		marker : nil,
		objectId: objectId,
	}
}

func (sc *sqlConfigure[T])Load() (T, error) {
	data_gen := new(T)
	data := *data_gen

	rows, err := sc.c.QueryContext(context.Background(), constant.SELECT_OBJECT_CONFIG_QUERY, sc.objectId)
	if err != nil {
		sc.levelLogger.Error("exec failed : ",err.Error())
		return data, err
	}

	key := ""
	value := ""

	for rows.Next() {
		scanErr := rows.Scan(&key, &value)
		if scanErr != nil {
			sc.levelLogger.Error("scan error : ", scanErr.Error())
			rows.Close()
			return data, scanErr
		}

		SetErr := data.Set(key, value)
		if SetErr != nil {
			sc.levelLogger.Warn("skip init this key: ", key)
		}
	}

	return data, nil
}


