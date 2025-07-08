package biz

import (
	"context"
	"im-server/internal/biz/do"
	"im-server/internal/pkg/infra/cache"
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo             GreeterRepo
	distributedCache cache.DistributedCache
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, distributedCache cache.DistributedCache) *GreeterUsecase {
	return &GreeterUsecase{
		repo:             repo,
		distributedCache: distributedCache,
	}
}

func (uc GreeterUsecase) ListAllGreeter(ctx context.Context) ([]*do.Greeter, error) {
	return uc.repo.ListAll(ctx)
}
