package biz

import (
	"context"
	"im-server/internal/biz/do"
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo) *GreeterUsecase {
	return &GreeterUsecase{repo: repo}
}

func (uc GreeterUsecase) ListAllGreeter(ctx context.Context) ([]*do.Greeter, error) {
	return uc.repo.ListAll(ctx)
}
