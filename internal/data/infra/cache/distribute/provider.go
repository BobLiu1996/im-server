package distribute

import (
	"github.com/google/wire"
	"im-server/internal/biz/do"
	"im-server/internal/pkg/infra/cache"
)

var ProviderSet = wire.NewSet(
	NewRedisDistributeCacheService,
	wire.Bind(
		new(cache.DistributedCache),
		new(*RedisDistributeCacheService),
	),
)

var GreeterProviderSet = wire.NewSet(
	ProvideGreeterCache,
)

func ProvideGreeterCache() cache.DistributedCacheType[*do.Greeter] {
	return NewRedisDistributeCacheType[*do.Greeter]()
}
