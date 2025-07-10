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
)

type (
	RedisData struct {
		// 实际业务数据
		Data any
		// 过期时间点
		ExpireTime time.Time
	}

	RedisDistributeCacheType[T any] struct {
		client redis.Cmdable
	}
)

func NewRedisDistributeCacheType[T any](data *data.Data) *RedisDistributeCacheType[T] {
	return &RedisDistributeCacheType[T]{
		client: data.Redis(),
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

func (r *RedisDistributeCacheType[T]) SetWithLogicalExpire(ctx context.Context, key string, value any, timeout time.Duration) error {
	redisData := &RedisData{
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
	// 缓存的数据为空字符串，直接返回nil
	if cachedValue == "" {
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
	// 缓存的数据为空字符串，直接返回nil
	if cachedValue == "" {
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
				if err = r.SetWithTTL(ctx, key, EmptyValue, CacheNullTTL); err != nil {
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
	if cachedValue == "" {
		return nil, nil
	}
	cachedList, err := GetResultList[T](cachedValue)
	if err != nil {
		return nil, err
	}
	return cachedList, nil
}

func (r *RedisDistributeCacheType[T]) QueryWithPassThroughListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpireWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpireList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithLogicalExpireListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithMutex(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithMutexWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (T, error), timeout time.Duration) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithMutexList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheType[T]) QueryWithMutexListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]T, error), timeout time.Duration) ([]T, error) {
	//TODO implement me
	panic("implement me")
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
	return fmt.Sprintf("%s_%s", keyPrefix, key)
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
