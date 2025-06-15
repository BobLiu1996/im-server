package service

import (
	"context"
	plog "im-server/pkg/log"
)

type CronService struct {
}

func NewCronService() *CronService {
	return &CronService{}
}

// InitService 初始化服务
func (cs *CronService) InitService() error {
	return nil
}

// a im-server that use cron task
func (cs *CronService) SyncStatus(ctx context.Context) error {
	plog.Info(ctx, "cron task sybcStatus is running")
	return nil
}
