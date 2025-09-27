package types

import (
	"agent_common/pkg/util/types"
	"io"
)

type InitData[PUSH any, CONN io.Closer, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] struct {
	SettingPath string
	CollectM map[string]CollectFn[PUSH, CONN]
	CronM map[string]CronFn[PUSH, FLAG, FLAGPTR]

	GetConnPoolFn GenCollectConnPoolFn[CONN]
	GetPusherFn GenPusherFn[PUSH]
}