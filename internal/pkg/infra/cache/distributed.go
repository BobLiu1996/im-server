package cache

import "time"

type DistributedCache interface {
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	Expire(key string, timeout time.Duration) bool
}
