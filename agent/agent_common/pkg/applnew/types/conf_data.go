package types

import (
	"strconv"
	"strings"
)

type ApplConfData struct {
	CollecbDbConfig struct {
		IP string
		Port int
		User string
		Password string	
		Dbname   string
	}
	
	Thread struct {
		CollectCount int
		PushCount    int
		CronCount    int
	}

	Queue struct {
		CollectSize  int
		CronSize     int
		PushSize     int
	}
}

func(aac *ApplConfData)Set(key string, value string) error {
	if strings.Index(key, "application_") == 0 {
		switch key {
		case "collectDBIP":
			aac.CollecbDbConfig.IP = value
		case "collectDBPort":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.CollecbDbConfig.Port = p
		case "collectDBUser":
			aac.CollecbDbConfig.User = value
		case "collectDBPasswd":
			aac.CollecbDbConfig.Password = value
		case "collectDBDbname":
			aac.CollecbDbConfig.Dbname = value
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
		case "threadCronCount":
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
		case "queueCronSize":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Queue.PushSize = p
		}
	}

	return nil
}