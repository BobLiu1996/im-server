package distribute

import (
	"context"
	"fmt"
	"im-server/internal/conf"
	"im-server/internal/data"
	"im-server/internal/pkg/infra/cache"
	"testing"
	"time"
)

const (
	keyPrefix = "test_distribute_cache_"
)

func InitRedisDistributeCacheService() cache.DistributedCache {
	config := &conf.Data{
		Mysql: &conf.Data_MySql{
			Driver: "mysql",
			Source: "root:root@tcp(192.168.5.134:3306)/test?charset=utf8mb4&parseTime=true&loc=Local",
		},
		Redis: &conf.Data_Redis{
			Addr:      "192.168.5.134:6379",
			Db:        0,
			Pool:      250,
			IsCluster: false,
		},
		Debug: true,
	}
	d, _, err := data.NewData(config)
	if err != nil {
		return nil
	}
	redisCache := NewRedisDistributeCacheService(d)
	return redisCache
}

func TestQueryWithPassThrough(t *testing.T) {
	redisCache := InitRedisDistributeCacheService()
	id := "test"
	d, err := redisCache.QueryWithPassThrough(context.Background(), keyPrefix, id, nil, 10*time.Second)
	if err != nil {
		_ = fmt.Errorf("cache query error: %v", err)
		return
	}
	fmt.Println(d)
}
