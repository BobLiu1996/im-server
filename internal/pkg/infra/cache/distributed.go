package cache

import (
	"context"
	"time"
)

// DistributedCache defines the interface for distributed caching.
type DistributedCache interface {
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, timeout time.Duration) error
	Expire(ctx context.Context, key string, timeout time.Duration) error
	SetWithLogicalExpire(ctx context.Context, key string, value any, timeout time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetObject(ctx context.Context, key string, target interface{}) (any, error)
	Delete(ctx context.Context, key string) error
	MultiGet(ctx context.Context, keys []string) (map[string]string, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	QueryWithPassThrough(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (any, error), timeout time.Duration) (any, error)
	QueryWithPassThroughWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (any, error), timeout time.Duration) (any, error)
	QueryWithPassThroughList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]any, error), timeout time.Duration) ([]any, error)
	QueryWithPassThroughListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]any, error), timeout time.Duration) ([]any, error)
	QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (any, error), timeout time.Duration) (any, error)
	QueryWithLogicalExpireWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (any, error), timeout time.Duration) (any, error)
	QueryWithLogicalExpireList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]any, error), timeout time.Duration) ([]any, error)
	QueryWithLogicalExpireListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]any, error), timeout time.Duration) (any, error)
	QueryWithMutex(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (any, error), timeout time.Duration) (any, error)
	QueryWithMutexWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (any, error), timeout time.Duration) (any, error)
	QueryWithMutexList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]any, error), timeout time.Duration) ([]any, error)
	QueryWithMutexListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]any, error), timeout time.Duration) (any, error)
}

type DistributedCacheType[T any] interface {
	Set(ctx context.Context, key string, value any) error
	QueryWithPassThroughList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, T) ([]T, error), timeout time.Duration) ([]T, error)
}
