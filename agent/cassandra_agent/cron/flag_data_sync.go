package cron

import (
	"agent_common/pkg/applnew/logger"
	"cassandra_agent/cache"
	"cassandra_agent/types"
	"context"
)

func cronFlagSync(ctx context.Context, flag *types.FlagData, log logger.LevelLogger) (*types.PushData, error){
	newTTL := flag.DataTTL.Load()
	oldTTL := cache.FlagCache.DataTTL.Swap(newTTL)

	if newTTL != oldTTL {
		log.Info("convert data TTL old :", oldTTL, " new :", newTTL)
	}

	return nil, nil
}