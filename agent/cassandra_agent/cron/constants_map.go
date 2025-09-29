package cron

import (
	appltype "agent_common/pkg/applnew/types"
	"cassandra_agent/types"
)


var CronMappsing  map[string]appltype.CronFn[types.PushData, types.FlagData, *types.FlagData] = map[string]appltype.CronFn[types.PushData, types.FlagData, *types.FlagData]{
	
}