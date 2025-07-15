//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"im-server/internal/biz"
	"im-server/internal/conf"
	"im-server/internal/data"
	rediscache "im-server/internal/data/infra/cache/distribute"
	localcache "im-server/internal/data/infra/cache/local"
	rocketmq "im-server/internal/data/infra/mq"

	redislock "im-server/internal/data/infra/lock/redis"
	"im-server/internal/server"
	"im-server/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.AppConfig, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, redislock.ProviderSet, rediscache.AutoWireDistributedCacheProviderSet, localcache.AutoWireLocalCacheProviderSet, rocketmq.ProviderSet, newApp))
}
