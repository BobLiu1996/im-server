package service

import (
	"github.com/google/wire"
	swire "im-server/internal/server/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewCronService,
	NewGreeterService,
	wire.Bind(new(swire.CronService), new(*CronService)),
)
