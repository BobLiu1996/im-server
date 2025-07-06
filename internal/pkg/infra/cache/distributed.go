package cache

import (
	"context"
	"time"
)

// DistributedCache defines the interface for distributed caching.
type DistributedCache[T any] interface {
	Set(key string, value T) error
	SetWithTTL(key string, value T, timeout time.Duration) error
	Expire(key string, timeout time.Duration) error
	SetWithLogicalExpire(key string, value T, timeout time.Duration) error
	Get(key string) (T, error)
	Delete(key string) error
	MultiGet(keys []string) (map[string]T, error)
	Keys(pattern string) ([]string, error)
	QueryWithPassThrough(keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error)
	QueryWithPassThroughWithoutArgs(keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error)
	QueryWithPassThroughList(keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithPassThroughListWithoutArgs(keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithLogicalExpire(keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error)
	QueryWithLogicalExpireWithoutArgs(keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error)
	QueryWithLogicalExpireList(keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithLogicalExpireListWithoutArgs(keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithMutex(keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error)
	QueryWithMutexWithoutArgs(keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error)
	QueryWithMutexList(keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error)
	QueryWithMutexListWithoutArgs(keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error)
}
