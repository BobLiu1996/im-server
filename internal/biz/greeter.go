package biz

import (
	"context"
	"fmt"
	"im-server/internal/biz/do"
	"im-server/internal/pkg/infra/cache"
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo                   GreeterRepo
	greeterDistributeCache cache.DistributedCacheType[*do.Greeter]
	userDistributeCache    cache.DistributedCacheType[*do.User]
	greeterLocalCache      cache.LocalCache[string, *do.Greeter]
	userLocalCache         cache.LocalCache[string, *do.User]
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, greeterDistributeCache cache.DistributedCacheType[*do.Greeter], userDistributeCache cache.DistributedCacheType[*do.User], greeterLocalCache cache.LocalCache[string, *do.Greeter], userLocalCache cache.LocalCache[string, *do.User]) *GreeterUsecase {
	return &GreeterUsecase{
		repo:                   repo,
		greeterDistributeCache: greeterDistributeCache,
		userDistributeCache:    userDistributeCache,
		greeterLocalCache:      greeterLocalCache,
		userLocalCache:         userLocalCache,
	}
}

func (uc GreeterUsecase) ListAllGreeter(ctx context.Context) ([]*do.Greeter, error) {
	// 获取Greeter类型的分布式缓存列表
	//greeterList, err := uc.greeterDistributeCache.QueryWithPassThroughList(ctx, "id", "", nil, 10*time.Second)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(greeterList)
	//// 获取User类型的分布式缓存列表
	//userList, err := uc.userDistributeCache.QueryWithPassThroughList(ctx, "id", "", nil, 10*time.Second)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(userList)
	// 获取Greeter类型的本地缓存
	present, err := uc.greeterLocalCache.GetIfPresent("test_local")
	if err != nil {
		return nil, err
	}
	fmt.Println(present)
	return uc.repo.ListAll(ctx)
}
