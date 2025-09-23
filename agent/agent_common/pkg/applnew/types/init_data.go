package types

import "agent_common/pkg/util/types"

type InitData[PUSH any, FLAG any, FLAGPTR types.GetterKeysetterInter[FLAG]] struct {
	SettingPath string
	CollectM map[string]CollectFn[PUSH]
	CronM map[string]CronFn[PUSH, FLAG, FLAGPTR]
	DataPusher DataPusher[PUSH]
}