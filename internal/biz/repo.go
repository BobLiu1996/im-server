package biz

import (
	"context"
	"im-server/internal/biz/do"
)

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	ListAll(context.Context) ([]*do.Greeter, error)
}
