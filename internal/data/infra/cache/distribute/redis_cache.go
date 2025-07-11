package distribute

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"im-server/internal/data"
	"im-server/internal/pkg/infra/lock"
	typeconversion "im-server/pkg/conversion"
	plog "im-server/pkg/log"
	"reflect"
	"strings"
	"time"
)

const (
	CacheNullTTL            = 60 * time.Second
	EmptyValue              = ""
	EmptyListValue          = "[]"
	LockSuffix              = "_lock"
	ThreadSleepMilliseconds = 50 * time.Millisecond
	LockExpiry              = 8 * time.Second
)

type (
	RedisData[T any] struct {
		// 实际业务数据
		Data T
		// 过期时间点
		ExpireTime time.Time
	}

	RedisDataList[T any] struct {
		// 实际业务数据
		Data []T
		// 过期时间点
		ExpireTime time.Time
	}

	RedisDistributeCacheType[T any] struct {
		client          redis.Cmdable
		distributedLock lock.DistributedLock
	}
)

func NewRedisDistributeCacheType[T any](data *data.Data, distributeLock lock.DistributedLock) *RedisDistributeCacheType[T] {
	return &RedisDistributeCacheType[T]{
		client:          data.Redis(),
		distributedLock: distributeLock,
	}
}

func (r *RedisDistributeCacheType[T]) Set(ctx context.Context, key string, value any) error {
	return r.client.Set(ctx, key, getValue(value), 0).Err()
}

func (r *RedisDistributeCacheType[T]) SetWithTTL(ctx context.Context, key string, value any, timeout time.Duration) error {
	return r.client.Set(ctx, key, getValue(value), timeout).Err()
}

func (r *RedisDistributeCacheType[T]) Expire(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

func (r *RedisDistributeCacheType[T]) SetWithLogicalExpire(ctx context.Context, key string, value T, timeout time.Duration) error {
	redisData := &RedisData[T]{
		Data:       value,
		ExpireTime: time.Now().Add(timeout),
	}
	return r.client.Set(ctx, key, getValue(redisData), 0).Err()
}

func (r *RedisDistributeCacheType[T]) SetWithLogicalExpireList(ctx context.Context, key string, value []T, timeout time.Duration) error {
	redisData := &RedisDataList[T]{
		Data:       value,
		ExpireTime: time.Now().Add(timeout),
	}
	return r.client.Set(ctx, key, getValue(redisData), 0).Err()
}

func (r *RedisDistributeCacheType[T]) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisDistributeCacheType[T]) GetObject(ctx context.Context, key string) (T, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return *new(T), nil
	} else if err != nil {
		return *new(T), err
	}
	result, err := GetResult[T](val)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *RedisDistributeCacheType[T]) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisDistributeCacheType[T]) MultiGet(ctx context.Context, keys []string) (map[string]string, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for i, key := range keys {
		if val, ok := vals[i].(string); ok {
			result[key] = val
		} else {
			result[key] = ""
		}
	}
	return result, nil
}

func (r *RedisDistributeCacheType[T]) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func (r *RedisDistributeCacheType[T]) QueryWithPassThrough(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	key := getKey(keyPrefix, id)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，查询数据库
			rVal, err := dbFallback(ctx, id)
			if err != nil {
				return Zero[T](), err
			}
			// 数据库中也不存在
			if IsEmpty(rVal) {
				// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
				if err = r.SetWithTTL(ctx, key, EmptyValue, CacheNullTTL); err != nil {
					return Zero[T](), err
				}
				return Zero[T](), nil
			}
			// 缓存数据
			if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
				// 如果缓存失败，直接返回数据库查询结果,同时返回错误
				return rVal, err
			}
			return rVal, nil
		}
		return Zero[T](), err
	}
	// 缓存的数据为空值，直接返回nil
	if cachedValue == EmptyValue {
		return Zero[T](), nil
	}
	result, err := GetResult[T](cachedValue)
	if err != nil {
		return Zero[T](), err
	}
	return result, nil
}

func (r *RedisDistributeCacheType[T]) QueryWithPassThroughWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	key := getKeyWithoutID(keyPrefix)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，查询数据库
			rVal, err := dbFallback(ctx)
			if err != nil {
				return Zero[T](), err
			}
			// 数据库中也不存在
			if IsEmpty(rVal) {
				// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
				if err = r.SetWithTTL(ctx, key, EmptyValue, CacheNullTTL); err != nil {
					return Zero[T](), err
				}
				return Zero[T](), nil
			}
			// 缓存数据
			if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
				// 如果缓存失败，直接返回数据库查询结果,同时返回错误
				return rVal, err
			}
			return rVal, nil
		}
		return Zero[T](), err
	}
	// 缓存的数据为空值，直接返回nil
	if cachedValue == EmptyValue {
		return Zero[T](), nil
	}
	result, err := GetResult[T](cachedValue)
	if err != nil {
		return Zero[T](), err
	}
	return result, nil
}

func (r *RedisDistributeCacheType[T]) QueryWithPassThroughList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	key := getKey(keyPrefix, id)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，查询数据库
			rVals, err := dbFallback(ctx, id)
			if err != nil {
				return nil, err
			}
			// 数据库中也不存在
			if rVals == nil || len(rVals) == 0 {
				// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
				if err = r.SetWithTTL(ctx, key, EmptyListValue, CacheNullTTL); err != nil {
					return nil, err
				}
				return nil, nil
			}
			// 缓存数据
			if err = r.SetWithTTL(ctx, key, rVals, timeout); err != nil {
				// 如果缓存失败，直接返回数据库查询结果,同时返回错误
				return rVals, err
			}
			return rVals, nil
		}
		return nil, err
	}
	// 缓存的数据为空字符串，直接返回nil
	if cachedValue == EmptyListValue {
		return nil, nil
	}
	cachedList, err := GetResultList[T](cachedValue)
	if err != nil {
		return nil, err
	}
	return cachedList, nil
}

func (r *RedisDistributeCacheType[T]) QueryWithPassThroughListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	key := getKeyWithoutID(keyPrefix)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，查询数据库
			rVals, err := dbFallback(ctx)
			if err != nil {
				return nil, err
			}
			// 数据库中也不存在
			if rVals == nil || len(rVals) == 0 {
				// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
				if err = r.SetWithTTL(ctx, key, EmptyListValue, CacheNullTTL); err != nil {
					return nil, err
				}
				return nil, nil
			}
			// 缓存数据
			if err = r.SetWithTTL(ctx, key, rVals, timeout); err != nil {
				// 如果缓存失败，直接返回数据库查询结果,同时返回错误
				return rVals, err
			}
			return rVals, nil
		}
		return nil, err
	}
	// 缓存的数据为空值，直接返回nil
	if cachedValue == EmptyValue {
		return nil, nil
	}
	cachedList, err := GetResultList[T](cachedValue)
	if err != nil {
		return nil, err
	}
	return cachedList, nil
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	key := getKey(keyPrefix, id)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中,异步重构缓存
		r.buildCache(ctx, id, dbFallback, timeout, key)
		// 等待缓存构建
		time.Sleep(ThreadSleepMilliseconds)
		// 重试，直到构建成功
		return r.QueryWithLogicalExpire(ctx, keyPrefix, id, dbFallback, timeout)
	} else if err != nil {
		// 其他未知错误，如网络中断等
		return Zero[T](), err
	}
	// 命中了缓存数据
	// 反序列化缓存数据
	var redisData RedisData[T]
	if err := json.Unmarshal([]byte(cachedValue), &redisData); err != nil {
		// 反序列化错误，直接返回空数据（注意：返回缓存的业务数据的空值，即redisData.Data，而不是缓存数据）
		return Zero[T](), err
	}
	// 检查缓存中的业务数据是否为空
	if IsEmpty[T](redisData.Data) {
		// 检查空值标记是否过期
		if time.Now().After(redisData.ExpireTime) {
			// 触发异步重建验证
			r.buildCache(ctx, id, dbFallback, timeout, key)
		}
		return Zero[T](), nil
	}
	// 检查是否已经逻辑过期
	if time.Now().Before(redisData.ExpireTime) {
		// 未过期，直接返回解析后的数据
		result, err := GetResult[T](redisData.Data)
		if err != nil {
			return Zero[T](), err
		}
		return result, nil
	}
	// 已经过期，触发异步构建缓存的流程，防止阻塞主线程
	// 只尝试获取一次分布式锁，避免多线程同时重建缓存，获取锁失败，证明有其他线程正在重建该缓存，所以可以直接退出而无需重试（意义不大且极大地降低了并发度）
	r.buildCache(ctx, id, dbFallback, timeout, key)
	// 返回已经过期的数据（数据的最终一致性）
	result, err := GetResult[T](redisData.Data)
	if err != nil {
		return Zero[T](), err
	}
	return result, nil
}

// buildCache 异步重建缓存
func (r *RedisDistributeCacheType[T]) buildCache(ctx context.Context, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration, key string) {
	lockKey := getLockKey(key)
	executeBuild := func() {
		unlock, err := r.distributedLock.Lock(ctx, lockKey, LockExpiry)
		if err != nil {
			// 获取分布式锁失败，直接退出，证明有其他线程正在重建该缓存
			return
		}
		defer unlock(ctx)
		// 获取锁成功，Double Check，再次检查缓存
		redisDataStr, err := r.client.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			plog.Errorf(ctx, "Double check failed: %v", err)
			return
		}
		// 判断缓存中的数据
		var newData T
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，从db中获取数据
			newData, err = dbFallback(ctx, id)
			if err != nil {
				plog.Errorf(ctx, "Get data from db failed: %v", err)
				return
			}
		} else {
			var redisData RedisData[T]
			if json.Unmarshal([]byte(redisDataStr), &redisData) == nil && time.Now().After(redisData.ExpireTime) {
				// 缓存已过期：查询数据库
				plog.Infof(ctx, "Cache expired, rebuilding cache for key: %s", key)
				newData, err = dbFallback(ctx, id)
			}
		}
		// 更新缓存
		if IsEmpty[T](newData) {
			if err := r.SetWithLogicalExpire(ctx, key, Zero[T](), CacheNullTTL); err != nil {
				plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
				return
			}
		} else {
			if err := r.SetWithLogicalExpire(ctx, key, newData, timeout); err != nil {
				plog.Errorf(ctx, "Set logical expire failed: %v", err)
				return
			}
		}
	}
	//// 异步执行重建缓存
	//go executeBuild()
	// 同步重建
	executeBuild()
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpireWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	key := getKeyWithoutID(keyPrefix)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中,异步重构缓存
		r.buildCacheWithoutArgs(ctx, dbFallback, timeout, key)
		// 等待缓存构建
		time.Sleep(ThreadSleepMilliseconds)
		// 重试，直到构建成功
		return r.QueryWithLogicalExpireWithoutArgs(ctx, keyPrefix, dbFallback, timeout)
	} else if err != nil {
		// 其他未知错误，如网络中断等
		return Zero[T](), err
	}
	// 命中了缓存数据
	// 反序列化缓存数据
	var redisData RedisData[T]
	if err := json.Unmarshal([]byte(cachedValue), &redisData); err != nil {
		// 反序列化错误，直接返回空数据（注意：返回缓存的业务数据的空值，即redisData.Data，而不是缓存数据）
		return Zero[T](), err
	}
	// 检查缓存中的业务数据是否为空
	if IsEmpty[T](redisData.Data) {
		// 检查空值标记是否过期
		if time.Now().After(redisData.ExpireTime) {
			// 触发异步重建验证
			r.buildCacheWithoutArgs(ctx, dbFallback, timeout, key)
		}
		return Zero[T](), nil
	}
	// 检查是否已经逻辑过期
	if time.Now().Before(redisData.ExpireTime) {
		// 未过期，直接返回解析后的数据
		result, err := GetResult[T](redisData.Data)
		if err != nil {
			return Zero[T](), err
		}
		return result, nil
	}
	// 已经过期，触发异步构建缓存的流程，防止阻塞主线程
	// 只尝试获取一次分布式锁，避免多线程同时重建缓存，获取锁失败，证明有其他线程正在重建该缓存，所以可以直接退出而无需重试（意义不大且极大地降低了并发度）
	r.buildCacheWithoutArgs(ctx, dbFallback, timeout, key)
	// 返回已经过期的数据（数据的最终一致性）
	result, err := GetResult[T](redisData.Data)
	if err != nil {
		return Zero[T](), err
	}
	return result, nil
}

// buildCache 异步重建缓存
func (r *RedisDistributeCacheType[T]) buildCacheWithoutArgs(ctx context.Context, dbFallback func(context.Context) (T, error), timeout time.Duration, key string) {
	lockKey := getLockKey(key)
	executeBuild := func() {
		unlock, err := r.distributedLock.Lock(ctx, lockKey, LockExpiry)
		if err != nil {
			// 获取分布式锁失败，直接退出，证明有其他线程正在重建该缓存
			return
		}
		defer unlock(ctx)
		// 获取锁成功，Double Check，再次检查缓存
		redisDataStr, err := r.client.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			plog.Errorf(ctx, "Double check failed: %v", err)
			return
		}
		// 判断缓存中的数据
		var newData T
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，从db中获取数据
			newData, err = dbFallback(ctx)
			if err != nil {
				plog.Errorf(ctx, "Get data from db failed: %v", err)
				return
			}
		} else {
			var redisData RedisData[T]
			if json.Unmarshal([]byte(redisDataStr), &redisData) == nil && time.Now().After(redisData.ExpireTime) {
				// 缓存已过期：查询数据库
				plog.Infof(ctx, "Cache expired, rebuilding cache for key: %s", key)
				newData, err = dbFallback(ctx)
			}
		}
		// 更新缓存
		if IsEmpty[T](newData) {
			if err := r.SetWithLogicalExpire(ctx, key, Zero[T](), CacheNullTTL); err != nil {
				plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
				return
			}
		} else {
			if err := r.SetWithLogicalExpire(ctx, key, newData, timeout); err != nil {
				plog.Errorf(ctx, "Set logical expire failed: %v", err)
				return
			}
		}
	}
	//// 异步执行重建缓存
	//go executeBuild()
	// 同步重建
	executeBuild()
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpireList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	key := getKey(keyPrefix, id)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中,异步重构缓存
		r.buildCacheList(ctx, id, dbFallback, timeout, key)
		// 等待缓存构建
		time.Sleep(ThreadSleepMilliseconds)
		// 重试，直到构建成功
		return r.QueryWithLogicalExpireList(ctx, keyPrefix, id, dbFallback, timeout)
	} else if err != nil {
		// 其他未知错误，如网络中断等
		return nil, err
	}
	// 命中了缓存数据
	// 反序列化缓存数据
	var redisData RedisDataList[T]
	if err := json.Unmarshal([]byte(cachedValue), &redisData); err != nil {
		// 反序列化错误，直接返回空数据（注意：返回缓存的业务数据的空值，即redisData.Data，而不是缓存数据）
		return nil, err
	}
	// 检查缓存中的业务数据是否为空
	if redisData.Data == nil || len(redisData.Data) == 0 {
		// 检查空值标记是否过期
		if time.Now().After(redisData.ExpireTime) {
			// 触发异步重建验证
			r.buildCacheList(ctx, id, dbFallback, timeout, key)
		}
		return nil, nil
	}
	// 检查是否已经逻辑过期
	if time.Now().Before(redisData.ExpireTime) {
		// 未过期，直接返回解析后的数据
		result, err := GetResultList[T](redisData.Data)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	// 已经过期，触发异步构建缓存的流程，防止阻塞主线程
	// 只尝试获取一次分布式锁，避免多线程同时重建缓存，获取锁失败，证明有其他线程正在重建该缓存，所以可以直接退出而无需重试（意义不大且极大地降低了并发度）
	r.buildCacheList(ctx, id, dbFallback, timeout, key)
	// 返回已经过期的数据（数据的最终一致性）
	result, err := GetResultList[T](redisData.Data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RedisDistributeCacheType[T]) buildCacheList(ctx context.Context, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration, key string) {
	lockKey := getLockKey(key)
	executeBuild := func() {
		unlock, err := r.distributedLock.Lock(ctx, lockKey, LockExpiry)
		if err != nil {
			// 获取分布式锁失败，直接退出，证明有其他线程正在重建该缓存
			return
		}
		defer unlock(ctx)
		// 获取锁成功，Double Check，再次检查缓存
		redisDataStr, err := r.client.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			plog.Errorf(ctx, "Double check failed: %v", err)
			return
		}
		// 判断缓存中的数据
		var newDataList []T
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，从db中获取数据
			newDataList, err = dbFallback(ctx, id)
			if err != nil {
				plog.Errorf(ctx, "Get data from db failed: %v", err)
				return
			}
		} else {
			var redisDataList RedisDataList[T]
			if json.Unmarshal([]byte(redisDataStr), &redisDataList) == nil && time.Now().After(redisDataList.ExpireTime) {
				// 缓存已过期：查询数据库
				plog.Infof(ctx, "Cache expired, rebuilding cache for key: %s", key)
				newDataList, err = dbFallback(ctx, id)
			}
		}
		// 更新缓存
		if newDataList == nil || len(newDataList) == 0 {
			if err := r.SetWithLogicalExpireList(ctx, key, nil, CacheNullTTL); err != nil {
				plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
				return
			}
		} else {
			if err := r.SetWithLogicalExpireList(ctx, key, newDataList, timeout); err != nil {
				plog.Errorf(ctx, "Set logical expire failed: %v", err)
				return
			}
		}
	}
	//// 异步执行重建缓存
	//go executeBuild()
	// 同步重建
	executeBuild()
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpireListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	key := getKeyWithoutID(keyPrefix)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中,异步重构缓存
		r.buildCacheListWithoutArgs(ctx, dbFallback, timeout, key)
		// 等待缓存构建
		time.Sleep(ThreadSleepMilliseconds)
		// 重试，直到构建成功
		return r.QueryWithLogicalExpireListWithoutArgs(ctx, keyPrefix, dbFallback, timeout)
	} else if err != nil {
		// 其他未知错误，如网络中断等
		return nil, err
	}
	// 命中了缓存数据
	// 反序列化缓存数据
	var redisData RedisDataList[T]
	if err := json.Unmarshal([]byte(cachedValue), &redisData); err != nil {
		// 反序列化错误，直接返回空数据（注意：返回缓存的业务数据的空值，即redisData.Data，而不是缓存数据）
		return nil, err
	}
	// 检查缓存中的业务数据是否为空
	if redisData.Data == nil || len(redisData.Data) == 0 {
		// 检查空值标记是否过期
		if time.Now().After(redisData.ExpireTime) {
			// 触发异步重建验证
			r.buildCacheListWithoutArgs(ctx, dbFallback, timeout, key)
		}
		return nil, nil
	}
	// 检查是否已经逻辑过期
	if time.Now().Before(redisData.ExpireTime) {
		// 未过期，直接返回解析后的数据
		result, err := GetResultList[T](redisData.Data)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	// 已经过期，触发异步构建缓存的流程，防止阻塞主线程
	// 只尝试获取一次分布式锁，避免多线程同时重建缓存，获取锁失败，证明有其他线程正在重建该缓存，所以可以直接退出而无需重试（意义不大且极大地降低了并发度）
	r.buildCacheListWithoutArgs(ctx, dbFallback, timeout, key)
	// 返回已经过期的数据（数据的最终一致性）
	result, err := GetResultList[T](redisData.Data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RedisDistributeCacheType[T]) buildCacheListWithoutArgs(ctx context.Context, dbFallback func(context.Context) ([]T, error), timeout time.Duration, key string) {
	lockKey := getLockKey(key)
	executeBuild := func() {
		unlock, err := r.distributedLock.Lock(ctx, lockKey, LockExpiry)
		if err != nil {
			// 获取分布式锁失败，直接退出，证明有其他线程正在重建该缓存
			return
		}
		defer unlock(ctx)
		// 获取锁成功，Double Check，再次检查缓存
		redisDataStr, err := r.client.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			plog.Errorf(ctx, "Double check failed: %v", err)
			return
		}
		// 判断缓存中的数据
		var newDataList []T
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，从db中获取数据
			newDataList, err = dbFallback(ctx)
			if err != nil {
				plog.Errorf(ctx, "Get data from db failed: %v", err)
				return
			}
		} else {
			var redisDataList RedisDataList[T]
			if json.Unmarshal([]byte(redisDataStr), &redisDataList) == nil && time.Now().After(redisDataList.ExpireTime) {
				// 缓存已过期：查询数据库
				plog.Infof(ctx, "Cache expired, rebuilding cache for key: %s", key)
				newDataList, err = dbFallback(ctx)
			}
		}
		// 更新缓存
		if newDataList == nil || len(newDataList) == 0 {
			if err := r.SetWithLogicalExpireList(ctx, key, nil, CacheNullTTL); err != nil {
				plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
				return
			}
		} else {
			if err := r.SetWithLogicalExpireList(ctx, key, newDataList, timeout); err != nil {
				plog.Errorf(ctx, "Set logical expire failed: %v", err)
				return
			}
		}
	}
	//// 异步执行重建缓存
	//go executeBuild()
	// 同步重建
	executeBuild()
}

func (r *RedisDistributeCacheType[T]) QueryWithMutex(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	key := getKey(keyPrefix, id)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中，尝试获取分布式锁
			unlock, err := r.distributedLock.Lock(ctx, getLockKey(key), LockExpiry)
			if err != nil {
				// 获取分布式锁失败，重试
				time.Sleep(ThreadSleepMilliseconds)
				return r.QueryWithMutex(ctx, keyPrefix, id, dbFallback, timeout)
			}
			defer unlock(ctx)
			// 再次检查缓存
			cachedValue, err = r.client.Get(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				// 缓存仍然不存在，查询数据库
				rVal, err := dbFallback(ctx, id)
				if err != nil {
					return Zero[T](), err
				}
				// 数据库中也不存在
				if IsEmpty(rVal) {
					// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
					if err = r.SetWithTTL(ctx, key, EmptyValue, CacheNullTTL); err != nil {
						return Zero[T](), err
					}
					return Zero[T](), nil
				}
				// 缓存数据
				if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
					// 如果缓存失败，直接返回数据库查询结果,同时返回错误
					plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
					return rVal, err
				}
				return rVal, nil
			} else if err != nil {
				return Zero[T](), err
			}
		} else {
			return Zero[T](), err // 其他错误，直接返回
		}
	}
	if cachedValue == EmptyValue {
		return Zero[T](), nil // 缓存的数据为空值，直接返回空值
	}
	result, err := GetResult[T](cachedValue)
	if err != nil {
		return Zero[T](), err // 反序列化错误，返回空值
	}
	return result, nil // 返回解析后的数据
}

func (r *RedisDistributeCacheType[T]) QueryWithMutexWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	key := getKeyWithoutID(keyPrefix)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中，尝试获取分布式锁
			unlock, err := r.distributedLock.Lock(ctx, getLockKey(key), LockExpiry)
			if err != nil {
				// 获取分布式锁失败，重试
				time.Sleep(ThreadSleepMilliseconds)
				return r.QueryWithMutexWithoutArgs(ctx, keyPrefix, dbFallback, timeout)
			}
			defer unlock(ctx)
			// 再次检查缓存
			cachedValue, err = r.client.Get(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				// 缓存仍然不存在，查询数据库
				rVal, err := dbFallback(ctx)
				if err != nil {
					return Zero[T](), err
				}
				// 数据库中也不存在
				if IsEmpty(rVal) {
					// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
					if err = r.SetWithTTL(ctx, key, EmptyValue, CacheNullTTL); err != nil {
						return Zero[T](), err
					}
					return Zero[T](), nil
				}
				// 缓存数据
				if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
					// 如果缓存失败，直接返回数据库查询结果,同时返回错误
					plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
					return rVal, err
				}
				return rVal, nil
			} else if err != nil {
				return Zero[T](), err
			}
		} else {
			return Zero[T](), err // 其他错误，直接返回
		}
	}
	if cachedValue == EmptyValue {
		return Zero[T](), nil // 缓存的数据为空值，直接返回空值
	}
	result, err := GetResult[T](cachedValue)
	if err != nil {
		return Zero[T](), err // 反序列化错误，返回空值
	}
	return result, nil // 返回解析后的数据
}

func (r *RedisDistributeCacheType[T]) QueryWithMutexList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	key := getKey(keyPrefix, id)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中，尝试获取分布式锁
			unlock, err := r.distributedLock.Lock(ctx, getLockKey(key), LockExpiry)
			if err != nil {
				// 获取分布式锁失败，重试
				time.Sleep(ThreadSleepMilliseconds)
				return r.QueryWithMutexList(ctx, keyPrefix, id, dbFallback, timeout)
			}
			defer unlock(ctx)
			// 再次检查缓存
			cachedValue, err = r.client.Get(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				// 缓存仍然不存在，查询数据库
				rVal, err := dbFallback(ctx, id)
				if err != nil {
					return nil, err
				}
				// 数据库中也不存在
				if rVal == nil || len(rVal) == 0 {
					// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
					if err = r.SetWithTTL(ctx, key, EmptyListValue, CacheNullTTL); err != nil {
						return nil, err
					}
					return nil, nil
				}
				// 缓存数据
				if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
					// 如果缓存失败，直接返回数据库查询结果,同时返回错误
					plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
					return rVal, err
				}
				return rVal, nil
			} else if err != nil {
				return nil, err
			}
		} else {
			return nil, err // 其他错误，直接返回
		}
	}
	if cachedValue == EmptyListValue {
		return nil, nil // 缓存的数据为空值，直接返回空值
	}
	result, err := GetResultList[T](cachedValue)
	if err != nil {
		return nil, err // 反序列化错误，返回空值
	}
	return result, nil // 返回解析后的数据
}

func (r *RedisDistributeCacheType[T]) QueryWithMutexListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	key := getKeyWithoutID(keyPrefix)
	// 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存未命中，尝试获取分布式锁
			unlock, err := r.distributedLock.Lock(ctx, getLockKey(key), LockExpiry)
			if err != nil {
				// 获取分布式锁失败，重试
				time.Sleep(ThreadSleepMilliseconds)
				return r.QueryWithMutexListWithoutArgs(ctx, keyPrefix, dbFallback, timeout)
			}
			defer unlock(ctx)
			// 再次检查缓存
			cachedValue, err = r.client.Get(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				// 缓存仍然不存在，查询数据库
				rVal, err := dbFallback(ctx)
				if err != nil {
					return nil, err
				}
				// 数据库中也不存在
				if rVal == nil || len(rVal) == 0 {
					// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
					if err = r.SetWithTTL(ctx, key, EmptyListValue, CacheNullTTL); err != nil {
						return nil, err
					}
					return nil, nil
				}
				// 缓存数据
				if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
					// 如果缓存失败，直接返回数据库查询结果,同时返回错误
					plog.Errorf(ctx, "Set empty logical expire failed: %v", err)
					return rVal, err
				}
				return rVal, nil
			} else if err != nil {
				return nil, err
			}
		} else {
			return nil, err // 其他错误，直接返回
		}
	}
	if cachedValue == EmptyListValue {
		return nil, nil // 缓存的数据为空值，直接返回空值
	}
	result, err := GetResultList[T](cachedValue)
	if err != nil {
		return nil, err // 反序列化错误，返回空值
	}
	return result, nil // 返回解析后的数据
}

func Zero[T any]() T {
	var zero T
	return zero
}

func IsEmpty[T any](v T) bool {
	// 特殊 case：T 是 interface{} 或 any，值可能是 nil
	if any(v) == nil {
		return true
	}
	val := reflect.ValueOf(v)
	// 如果是切片、map、chan、pointer等 nil-able 类型，判断是否为 nil
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return val.IsNil()
	}
	// 否则判断是否为零值（结构体/基本类型）
	return val.IsZero()
}

// getKey 获取带Id的缓存键
func getKey(keyPrefix string, id interface{}) string {
	return getKeyWithID(keyPrefix, id)
}

func getLockKey(key string) string {
	return key + LockSuffix
}

// getValue 获取要保存到缓存的值，可能是简单的类型，可能是对象类型，也可能是数组类型等
func getValue(obj any) string {
	if obj == nil {
		return ""
	}
	return toJSONString(obj)
}

// getKeyWithoutID 获取不带Id的缓存键
func getKeyWithoutID(keyPrefix string) string {
	return getKeyWithID(keyPrefix, nil)
}

// getKeyWithID 获取带有参数的缓存键
func getKeyWithID(keyPrefix string, id interface{}) string {
	if id == nil {
		return keyPrefix
	}
	tc := typeconversion.NewTypeConversion()
	var key string
	if tc.IsSimpleType(id) {
		key = fmt.Sprintf("%v", id)
	} else {
		jsonStr := toJSONString(id)
		hash := md5.Sum([]byte(jsonStr))
		key = hex.EncodeToString(hash[:])
	}
	if strings.TrimSpace(key) == "" {
		key = ""
	}
	return fmt.Sprintf("%s%s", keyPrefix, key)
}

// toJSONString 将对象转换为JSON字符串
func toJSONString(obj interface{}) string {
	s, ok := obj.(string)
	if ok {
		return s
	}
	d, err := json.Marshal(obj)
	if err != nil {
		plog.Errorf(context.Background(), "Failed to marshal object to JSON: %v", err)
		return ""
	}
	return string(d)
}

// GetResult 将json字符串转换成对象
func GetResult[T any](obj any) (T, error) {
	str := toJSONString(obj)
	var t T
	if err := json.Unmarshal([]byte(str), &t); err != nil {
		return *new(T), err // 返回T类型的零值
	}
	return t, nil
}

// GetResultList 将数组类型的json字符串转换成对象列表
func GetResultList[T any](obj any) ([]T, error) {
	if obj == nil {
		return nil, nil
	}
	str := toJSONString(obj)
	var t []T
	if err := json.Unmarshal([]byte(str), &t); err != nil {
		return nil, err
	}
	return t, nil
}
