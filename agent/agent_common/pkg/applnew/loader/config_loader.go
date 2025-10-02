package loader

import (
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/util/types"
	"context"
	"database/sql"
)

type Configure[
	CONF any, SYNC any, FLAG any,
	CONFPTR types.SetterInter[CONF], SYNCPTR types.SetterInter[SYNC], FLAGPTR types.GetterKeysetterInter[FLAG]] interface {
	LoadConfig() (CONFPTR, error)
	LoadSync() (SYNCPTR, error)
	LoadFlag() (FLAGPTR, error)
}

type sqlConfigure[
	CONF any, SYNC any, FLAG any,
	CONFPTR types.SetterInter[CONF], SYNCPTR types.SetterInter[SYNC], FLAGPTR types.GetterKeysetterInter[FLAG]] struct {
	c        *sql.DB
	objectId int
	otherLogger logger.LevelLogger

	flagHist FLAGPTR
}

func (s *sqlConfigure[CONF, SYNC, FLAG, CONFPTR, SYNCPTR, FLAGPTR]) LoadConfig() (CONFPTR, error) {
	data_gen := new(CONF)
	var data CONFPTR = data_gen

	rows, err := s.c.QueryContext(context.Background(), _SELECT_OBJECT_CONFIG_QUERY, s.objectId)
	if err != nil {
		s.otherLogger.Error("exec failed : ",err.Error())
		return data, err
	}

	key := ""
	value := ""

	for rows.Next() {
		scanErr := rows.Scan(&key, &value)
		if scanErr != nil {
			s.otherLogger.Error("scan error : ", scanErr.Error())
			rows.Close()
			return data, scanErr
		}

		SetErr := data.Set(key, value)
		if SetErr != nil {
			s.otherLogger.Warn("skip init this key: ", key)
		}
	}

	return data, nil
}

func (s *sqlConfigure[CONF, SYNC, FLAG, CONFPTR, SYNCPTR, FLAGPTR]) LoadFlag() (FLAGPTR, error) {
	data_gen := new(FLAG)
	var data FLAGPTR = data_gen
	
	rows, err := s.c.QueryContext(context.Background(), _SELECT_OBJECT_CONFIG_QUERY, s.objectId)
	if err != nil {
		s.otherLogger.Error("exec failed : ",err.Error())
		return data, err
	}

	key := ""
	value := ""

	for rows.Next() {
		scanErr := rows.Scan(&key, &value)
		if scanErr != nil {
			s.otherLogger.Error("scan error : ", scanErr.Error())
			rows.Close()
			return data, scanErr
		}

		SetErr := data.Set(key, value)
		if SetErr != nil {
			s.otherLogger.Warn("skip init this key: ", key)
		}
	}

	return data, nil
}

func (s *sqlConfigure[CONF, SYNC, FLAG, CONFPTR, SYNCPTR, FLAGPTR]) LoadSync() (SYNCPTR, error) {
	data_gen := new(SYNC)
	var data SYNCPTR = data_gen

	rows, err := s.c.QueryContext(context.Background(), _SELECT_OBJECT_CONFIG_QUERY, s.objectId)
	if err != nil {
		s.otherLogger.Error("exec failed : ",err.Error())
		return data, err
	}

	key := ""
	value := ""

	for rows.Next() {
		scanErr := rows.Scan(&key, &value)
		if scanErr != nil {
			s.otherLogger.Error("scan error : ", scanErr.Error())
			rows.Close()
			return data, scanErr
		}

		SetErr := data.Set(key, value)
		if SetErr != nil {
			s.otherLogger.Warn("skip init this key: ", key)
		}
	}

	return data, nil
}

func NewSQLConfigure[
	CONF any, SYNC any, FLAG any,
	CONFPTR types.SetterInter[CONF], SYNCPTR types.SetterInter[SYNC], FLAGPTR types.GetterKeysetterInter[FLAG]](
	conn *sql.DB, objectId int, otherLogger logger.LevelLogger) Configure[CONF, SYNC, FLAG, CONFPTR, SYNCPTR, FLAGPTR] {
	return &sqlConfigure[CONF, SYNC, FLAG, CONFPTR, SYNCPTR, FLAGPTR]{
		c: conn,
		objectId:    objectId,
		otherLogger: otherLogger,
	}
}
