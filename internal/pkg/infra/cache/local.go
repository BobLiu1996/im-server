package cache

type LocalCache interface {
	Put(key interface{}, value interface{}) error
	GetIfPresent(key interface{}) (interface{}, error)
	Remove(key interface{}) error
}
