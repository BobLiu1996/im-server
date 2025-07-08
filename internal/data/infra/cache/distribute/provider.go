package distribute

import (
	"github.com/google/wire"
	"im-server/internal/pkg/infra/cache"
)

var ProviderSet = wire.NewSet(
	NewRedisDistributeCacheService,
	wire.Bind(
		new(cache.DistributedCache),
		new(*RedisDistributeCacheService),
	),
)
