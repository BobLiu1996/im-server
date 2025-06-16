package data

import (
	"context"
	"im-server/internal/biz"
	"im-server/internal/biz/do"
)

type mockGreeterRepo struct{}

func NewMockGreeterRepo() biz.GreeterRepo {
	return &mockGreeterRepo{}
}

func (m mockGreeterRepo) ListAll(ctx context.Context) ([]*do.Greeter, error) {
	//TODO implement me
	panic("implement me")
}
