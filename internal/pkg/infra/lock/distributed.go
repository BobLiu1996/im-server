package lock

import (
	"context"
	"time"
)

type Locker interface {
	// Lock acquires a lock with the given key and expiration.
	Lock(ctx context.Context, key string, timeout time.Duration) (unlock func(context.Context) error, err error)
}
