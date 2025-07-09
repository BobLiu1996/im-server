package biz

import (
	"context"
	"fmt"
	"im-server/internal/biz/do"
	"im-server/internal/pkg/infra/cache"
	"time"
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo                   GreeterRepo
	distributedCache       cache.DistributedCache
	greeterDistributeCache cache.DistributedCacheType[*do.Greeter]
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, distributedCache cache.DistributedCache, greeterDistributeCache cache.DistributedCacheType[*do.Greeter]) *GreeterUsecase {
	return &GreeterUsecase{
		repo:                   repo,
		distributedCache:       distributedCache,
		greeterDistributeCache: greeterDistributeCache,
	}
}

func (uc GreeterUsecase) ListAllGreeter(ctx context.Context) ([]*do.Greeter, error) {
	list, err := uc.greeterDistributeCache.QueryWithPassThroughList(ctx, "id", "", nil, 10*time.Second)
	if err != nil {
		return nil, err
	}
	fmt.Println(list)
	return uc.repo.ListAll(ctx)
}
