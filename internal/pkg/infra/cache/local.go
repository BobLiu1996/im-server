package cache

type LocalCache[K comparable, V any] interface {
	Put(key K, value V)
	GetIfPresent(key K) (V, error)
	Remove(key K)
}
