package types

import (
	"strconv"
	"strings"
	"sync"
)

type AppSyncData struct {
	Intervals map[string]int
	Cron      map[string]int
	Custom    map[string]string
	Stop     bool
	once sync.Once
}

func(asd *AppSyncData)Set(key string, value string) error {
	const intvlStartOff = len("interval")
	const cronStartOff = len("cron")

	asd.once.Do(func() {
		asd.Intervals = make(map[string]int)
	})

	if strings.ContainsAny(key, "interval") {
		p, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		
		asd.Intervals[key[intvlStartOff:]] = p
	} else if strings.ContainsAny(key, "cron") {
		p, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		
		asd.Cron[key[cronStartOff:]] = p
	}else if key == "stopAgnet" {
		if value == "Y" {
			asd.Stop = true
		} else {
			asd.Stop = false
		}
	} else {
		asd.Custom[key] = value
	}
	
	return nil
}