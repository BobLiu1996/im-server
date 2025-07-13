package local

import (
	"github.com/google/wire"
	"im-server/internal/biz/do"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/cache"
)

var AutoWireLocalCacheProviderSet = wire.NewSet(
	ProvideGreeterLocalCache,
	ProvideUserLocalCache,
)

// ProvideGreeterLocalCache 提供Greeter类型的本地缓存
// 该函数作用为接口绑定，相当于wire.Bind
func ProvideGreeterLocalCache(dataConfig *conf.Data) cache.LocalCache[string, *do.Greeter] {
	return NewLocalGoCache[string, *do.Greeter](dataConfig)
}

// ProvideUserLocalCache 提供User类型的本地缓存
func ProvideUserLocalCache(dataConfig *conf.Data) cache.LocalCache[string, *do.User] {
	return NewLocalGoCache[string, *do.User](dataConfig)
}
