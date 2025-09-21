package config

import "agent_common/pkg/util/types"

type DBConfig struct {
	Cassnadra struct {
		Addr string
		User string
		Password string
	}

	CollectDb struct {
		Addr string
		User string
		Password string
	}
}

func (cpc *DBConfig) CloneFromGetter(getter types.StrGetter) {


}