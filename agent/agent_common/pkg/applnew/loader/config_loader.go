package loader

import (
	"agent_common/pkg/applnew/logger"
	"agent_common/pkg/util/types"
	"context"
	"database/sql"
	"fmt"
)


type ConfigureUpdater[FLAG any,FLAGPTR types.GetterKeysetterInter[FLAG]] interface {
	UpdateFlag(FLAGPTR) error
}

type Configure[
	CONF any, SYNC any, FLAG any,
	CONFPTR types.SetterInter[CONF], SYNCPTR types.SetterInter[SYNC], FLAGPTR types.GetterKeysetterInter[FLAG]] interface {
	ConfigureUpdater[FLAG,FLAGPTR]
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

func (s *sqlConfigure[CONF, SYNC, FLAG, CONFPTR, SYNCPTR, FLAGPTR]) UpdateFlag(flag FLAGPTR) error {
	if flag == nil {
		s.otherLogger.Error("flag is nil pointer")
		return fmt.Errorf("configure, flag is nil pointer")
	}
	
	tx, txErr := s.c.Begin()
	if txErr != nil {
		s.otherLogger.Error("get trans error : ", txErr.Error())
		return txErr
	}

	for _, key := range flag.Keys() {
		val, getErr := flag.Get(key)
		if getErr != nil {
			tx.Rollback()
			s.otherLogger.Error("flag Key get Value error :", getErr.Error())
			return getErr
		}

		_, retErr := tx.Exec(_UPDATE_OBJECT_FLAG_QUERY, s.objectId, val)
		if retErr != nil {
			tx.Rollback()
			s.otherLogger.Error("flag update error :", retErr.Error())
			return retErr
		}
	}

	return tx.Commit()
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
