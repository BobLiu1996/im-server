package distribute

import (
	"context"
	"github.com/redis/go-redis/v9"
	"im-server/internal/data"
	distributed_cache "im-server/internal/pkg/infra/cache"
	"time"
)

type (
	RedisDistributeCacheService[T any] struct {
		client redis.Cmdable
	}
)

func NewRedisDistributeCacheService[T any](data *data.Data) distributed_cache.DistributedCache[T] {
	return &RedisDistributeCacheService[T]{client: data.Redis()}
}

func (r RedisDistributeCacheService[T]) Set(key string, value T) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) SetWithTTL(key string, value T, timeout time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) Expire(key string, timeout time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) SetWithLogicalExpire(key string, value T, timeout time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) Get(key string) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) Delete(key string) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) MultiGet(keys []string) (map[string]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) Keys(pattern string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithPassThrough(keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithPassThroughWithoutArgs(keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithPassThroughList(keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithPassThroughListWithoutArgs(keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithLogicalExpire(keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithLogicalExpireWithoutArgs(keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithLogicalExpireList(keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithLogicalExpireListWithoutArgs(keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithMutex(keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithMutexWithoutArgs(keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithMutexList(keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisDistributeCacheService[T]) QueryWithMutexListWithoutArgs(keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}
