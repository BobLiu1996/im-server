package redis

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
)

var (
	//go:embed script/unlock.lua
	luaUnlock string
	//go:embed script/refresh.lua
	luaRefresh string

	//go:embed script/lock.lua
	luaLock string

	ErrFailedToPreemptLock = errors.New("rlock: 抢锁失败")
	// ErrLockNotHold 一般是出现在你预期你本来持有锁，结果却没有持有锁的地方
	// 比如说当你尝试释放锁的时候，可能得到这个错误
	// 这一般意味着有人绕开了 rlock 的控制，直接操作了 Redis
	ErrLockNotHold = errors.New("rlock: 未持有锁")
)

type Locker struct {
	client redis.Cmdable
	g      singleflight.Group
	// valuer 用于生成值，将来可以考虑暴露出去允许用户自定义
	valuer func() string
}

func NewLocker(client redis.Cmdable) *Locker {
	return &Locker{
		client: client,
		valuer: func() string {
			return uuid.New().String()
		},
	}
}

func (c *Locker) Lock(ctx context.Context, key string, timeout time.Duration) (unlock func(context.Context) error, err error) {
	val := c.valuer()
	ok, err := c.client.SetNX(ctx, key, val, timeout).Result()
	if err != nil {
		// 网络问题，服务器问题，或者超时，都会走过来这里
		return nil, err
	}
	if !ok {
		// 已经有人加锁了，或者刚好和人一起加锁，但是自己竞争失败了
		return nil, ErrFailedToPreemptLock
	}
	return newLock(c.client, key, val, timeout).Unlock, nil
}

type Lock struct {
	client           redis.Cmdable
	key              string
	value            string
	expiration       time.Duration
	unlock           chan struct{}
	signalUnlockOnce sync.Once
}

func newLock(client redis.Cmdable, key string, value string, expiration time.Duration) *Lock {
	return &Lock{
		client:     client,
		key:        key,
		value:      value,
		expiration: expiration,
		unlock:     make(chan struct{}, 1),
	}
}

func (l *Lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	ticker := time.NewTicker(interval)
	// 刷新超时 channel
	ch := make(chan struct{}, 1)
	defer func() {
		ticker.Stop()
		close(ch)
	}()
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			// 超时这里，可以继续尝试
			if err == context.DeadlineExceeded {
				// 因为有两个可能的地方要写入数据，而 ch
				// 容量只有一个，所以如果写不进去就说明前一次调用超时了，并且还没被处理，
				// 与此同时计时器也触发了
				select {
				case ch <- struct{}{}:
				default:
				}
				continue
			}
			if err != nil {
				return err
			}
		case <-ch:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			// 超时这里，可以继续尝试
			if err == context.DeadlineExceeded {
				select {
				case ch <- struct{}{}:
				default:
				}
				continue
			}
			if err != nil {
				return err
			}
		case <-l.unlock:
			return nil
		}
	}
}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh,
		[]string{l.key}, l.value, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}

// Unlock 解锁
func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnlock, []string{l.key}, l.value).Int64()
	defer func() {
		// 避免重复解锁引起 panic
		l.signalUnlockOnce.Do(func() {
			l.unlock <- struct{}{}
			close(l.unlock)
		})
	}()
	if err == redis.Nil {
		return ErrLockNotHold
	}
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
	return nil
}
