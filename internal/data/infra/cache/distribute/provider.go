package distribute

import (
	"github.com/google/wire"
	icache "im-server/internal/pkg/infra/cache"
	ilock "im-server/internal/pkg/infra/lock"
)

var ProviderSet = wire.NewSet(NewRedisDistributeCacheService, wire.Bind(new(icache.DistributedCache[T any]), new(*RedisDistributeCacheService)))
