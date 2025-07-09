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
	"im-server/internal/pkg/infra/cache"
	typeconversion "im-server/pkg/conversion"
	plog "im-server/pkg/log"
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

var _ cache.DistributedCache = (*RedisDistributeCacheService)(nil)

type (
	RedisData struct {
		// 实际业务数据
		Data any
		// 过期时间点
		ExpireTime time.Time
	}

	RedisDistributeCacheService struct {
		client redis.Cmdable
	}
)

func NewRedisDistributeCacheService(data *data.Data) *RedisDistributeCacheService {
	return &RedisDistributeCacheService{client: data.Redis()}
}

func (r *RedisDistributeCacheService) Set(ctx context.Context, key string, value any) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisDistributeCacheService) SetWithTTL(ctx context.Context, key string, value any, timeout time.Duration) error {
	return r.client.Set(ctx, key, value, timeout).Err()
}

func (r *RedisDistributeCacheService) Expire(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

func (r *RedisDistributeCacheService) SetWithLogicalExpire(ctx context.Context, key string, value any, timeout time.Duration) error {
	redisData := &RedisData{
		Data:       value,
		ExpireTime: time.Now().Add(timeout),
	}
	return r.client.Set(ctx, key, redisData, 0).Err()
}

func (r *RedisDistributeCacheService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisDistributeCacheService) GetObject(ctx context.Context, key string, target interface{}) (any, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("key not found")
	} else if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(val), target); err != nil {
		return nil, err
	}
	return target, nil
}

func (r *RedisDistributeCacheService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisDistributeCacheService) MultiGet(ctx context.Context, keys []string) (map[string]string, error) {
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

func (r *RedisDistributeCacheService) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func (r *RedisDistributeCacheService) QueryWithPassThrough(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (any, error), timeout time.Duration) (any, error) {
	key := getKey(keyPrefix, id)
	// 2. 尝试从缓存获取
	cachedValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存不存在，查询数据库
			rVal, err := dbFallback(ctx, id)
			if err != nil {
				return nil, err
			}
			// 数据库中也不存在
			if rVal == nil {
				// 数据库为空，缓存空值（防止缓存穿透的关键步骤）
				if err = r.SetWithTTL(ctx, key, EmptyValue, CacheNullTTL); err != nil {
					return nil, err
				}
				return nil, nil
			}
			// 缓存数据
			if err = r.SetWithTTL(ctx, key, rVal, timeout); err != nil {
				return nil, err
			}
			return rVal, nil
		}
		return nil, err
	}
	// 缓存的数据为空字符串，直接返回nil
	if cachedValue == "" {
		return nil, nil
	}
	return cachedValue, nil
}

func (r *RedisDistributeCacheService) QueryWithPassThroughWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithPassThroughList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]any, error), timeout time.Duration) ([]any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithPassThroughListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]any, error), timeout time.Duration) ([]any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithLogicalExpire(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithLogicalExpireWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithLogicalExpireList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]any, error), timeout time.Duration) ([]any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithLogicalExpireListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithMutex(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) (any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithMutexWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) (any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithMutexList(ctx context.Context, keyPrefix string, id any, dbFallback func(context.Context, any) ([]any, error), timeout time.Duration) ([]any, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisDistributeCacheService) QueryWithMutexListWithoutArgs(ctx context.Context, keyPrefix string, dbFallback func(context.Context) ([]any, error), timeout time.Duration) (any, error) {
	//TODO implement me
	panic("implement me")
}

// getKey 获取缓存键
// 默认实现，调用带id参数的getKey
func getKey(keyPrefix string, id interface{}) string {
	return getKeyWithID(keyPrefix, id)
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

func getResult[T any](obj interface{}) T {
	str := toJSONString(obj)
	var t T
	err := json.Unmarshal([]byte(str), &t)
	if err != nil {
		plog.Errorf(context.Background(), "Failed to unmarshal JSON to type %T: %v", t, err)
		return *new(T) // 返回T类型的零值
	}
	return t
}
