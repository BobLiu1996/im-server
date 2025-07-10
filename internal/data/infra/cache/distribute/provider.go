package distribute

import (
	"github.com/google/wire"
	"im-server/internal/biz/do"
	"im-server/internal/data"
	"im-server/internal/pkg/infra/cache"
)

var AutoWireDistributedCacheProviderSet = wire.NewSet(
	ProvideGreeterDistributeCache,
	ProvideUserDistributeCache,
)

// ProvideGreeterDistributeCache 提供Greeter类型的分布式缓存
// 该函数作用为接口绑定，相当于wire.Bind
func ProvideGreeterDistributeCache(data *data.Data) cache.DistributedCacheType[*do.Greeter] {
	return NewRedisDistributeCacheType[*do.Greeter](data)
}

// ProvideUserDistributeCache 提供User类型的分布式缓存
func ProvideUserDistributeCache(data *data.Data) cache.DistributedCacheType[*do.User] {
	return NewRedisDistributeCacheType[*do.User](data)
}
