package appl

import (
	"strconv"
	"strings"
)

type AgentApplicationInitConfig struct {
	ObjectId  int
	LogConfig struct {
		Level string
		Path  string
	}

	ConfigDb struct {
		Drvier string
		Dsn    string
	}

	ConfigType string
}

type AgentApplicationConfig struct {
	collecbDbConfig struct {
		ip string
		port int
		user string
		password string	
		dbname   string
	}
	
	Thread struct {
		CollectCount int
		PushCount    int
	}

	Queue struct {
		CollectSize  int
		PushSize     int
	}

	intervals []struct {
		key string
		inetervalSec int

	}
	etc       map[string]string
}

func(aac *AgentApplicationConfig)Get(key string) (string, error) {

	return aac.etc[key], nil
}

func(aac *AgentApplicationConfig)Set(key string, value string) error {
	if strings.Index(key, "application_") == 0 {
		switch key {
		case "collectDBIP":
			aac.collecbDbConfig.ip = value
		case "collectDBPort":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.collecbDbConfig.port = p
		case "collectDBUser":
			aac.collecbDbConfig.user = value
		case "collectDBPasswd":
			aac.collecbDbConfig.password = value
		case "collectDBDbname":
			aac.collecbDbConfig.dbname = value
		case "threadCollectCount":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Thread.CollectCount = p
		case "threadPushCount":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Thread.PushCount = p
		case "queueCollectSize":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Queue.CollectSize = p
		case "queuePushSize":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Queue.PushSize = p
		}
	}

	if strings.Index(key, "interval_") == 0 {
		if strings.ContainsAny(key, "interval_") {
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}

			if p > 0 {
				aac.intervals = append(aac.intervals, struct{key string; inetervalSec int}{key, p})
			}
		} 
	}else {
		aac.etc[key] = value
	}
	return nil
}