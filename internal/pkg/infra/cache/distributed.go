package cache

import (
	"context"
	"time"
)

// DistributedCacheType defines the interface for a distributed cache with generic type T.
type DistributedCacheType[T any] interface {
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, timeout time.Duration) error
	Expire(ctx context.Context, key string, timeout time.Duration) error
	SetWithLogicalExpire(ctx context.Context, key string, value T, timeout time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetObject(ctx context.Context, key string) (T, error)
	Delete(ctx context.Context, key string) error
	MultiGet(ctx context.Context, keys []string) (map[string]string, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	QueryWithPassThrough(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error)
	QueryWithPassThroughWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error)
	QueryWithPassThroughList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithPassThroughListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error)
	QueryWithLogicalExpireWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error)
	QueryWithLogicalExpireList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithLogicalExpireListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithMutex(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error)
	QueryWithMutexWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error)
	QueryWithMutexList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithMutexListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error)
}
