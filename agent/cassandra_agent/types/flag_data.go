package types

import (
	"fmt"
	"strconv"
	"sync/atomic"
)

type FlagData struct {
	DataTTL atomic.Int64
}

type FlagDataKey string

const (
	DataTTLKey FlagDataKey = "DataTTLKey"
)

func (fd *FlagData)Set(key string, val string) (err error) {
	switch FlagDataKey(key) {
	case DataTTLKey:
		data,err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		fd.DataTTL.Store(int64(data))
	}
	return nil
}

func (fd *FlagData)Keys() []string {
	return []string{string(DataTTLKey)}
}

func (fd *FlagData)Get(key string) (string, error) {
	switch FlagDataKey(key) {
	case DataTTLKey:
		val := strconv.Itoa(int(fd.DataTTL.Load()))
		return val, nil
	}

	return "",fmt.Errorf("not exists key : %s", key)
}

