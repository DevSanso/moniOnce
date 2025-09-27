package types

import (

)

type FlagData struct {}

type FlagDataKey string



func (fd *FlagData)Set(key string, val string) error {
	return nil
}

func (fd *FlagData)Keys() []string {
	return []string{}
}

func (fd *FlagData)Get(key string) (string, error) {
	return "",nil
}