package types

import (
	"strconv"
	"strings"
)

type ApplConfData struct {
	Version float64
	
	CollecbDbConfig struct {
		Url string
	}

	DataPusherConfig struct {
		Url string
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
	if strings.Index(key, "application.") == 0 {
		switch key {
		case "collect.db.url":
			aac.CollecbDbConfig.Url = value
		case "thread.collectCount":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Thread.CollectCount = p
		case "thread.pushCount":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Thread.PushCount = p
		case "thread.cronCount":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Thread.PushCount = p
		case "queue.collectSize":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Queue.CollectSize = p
		case "queue.pushSize":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Queue.PushSize = p
		case "queue.cronSize":
			p, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			aac.Queue.PushSize = p
		case "dataPusher.conn.url":
			aac.DataPusherConfig.Url = value
		}
	}

	return nil
}