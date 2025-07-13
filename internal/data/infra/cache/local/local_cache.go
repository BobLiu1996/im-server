package local

import (
	"github.com/patrickmn/go-cache"
	"im-server/internal/conf"
	"im-server/internal/data/infra"
)

type LocalGoCache[K comparable, V any] struct {
	cache *cache.Cache
}

func NewLocalGoCache[K comparable, V any](dataConfig *conf.Data) *LocalGoCache[K, V] {
	return &LocalGoCache[K, V]{
		cache: cache.New(dataConfig.GetLocalCache().GetExpiration().AsDuration(), dataConfig.GetLocalCache().GetCleanupInterval().AsDuration()),
	}
}

func (l *LocalGoCache[K, V]) Put(key K, value V) {
	l.cache.Set(infra.GetValue(key), infra.GetValue(value), 0)
}

func (l *LocalGoCache[K, V]) GetIfPresent(key K) (V, error) {
	val, ok := l.cache.Get(infra.GetValue(key))
	if ok {
		result, err := infra.GetResult[V](val)
		if err != nil {
			return infra.Zero[V](), err
		}
		return result, nil
	}
	return infra.Zero[V](), nil
}

func (l *LocalGoCache[K, V]) Remove(key K) {
	l.cache.Delete(infra.GetValue(key))
}
