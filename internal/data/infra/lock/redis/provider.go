package redis

import (
	"github.com/google/wire"
	ilock "im-server/internal/pkg/infra/lock"
)

var ProviderSet = wire.NewSet(NewLocker, wire.Bind(new(ilock.DistributedLock), new(*Locker)))
