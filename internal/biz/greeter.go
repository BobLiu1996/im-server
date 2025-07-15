package biz

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"im-server/internal/biz/do"
	"im-server/internal/pkg/infra/cache"
	"im-server/internal/pkg/infra/model"
	"im-server/internal/pkg/infra/mq"
)

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo                   GreeterRepo
	greeterDistributeCache cache.DistributedCacheType[*do.Greeter]
	userDistributeCache    cache.DistributedCacheType[*do.User]
	greeterLocalCache      cache.LocalCache[string, *do.Greeter]
	userLocalCache         cache.LocalCache[string, *do.User]
	rmqSender              mq.MessageSender
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, greeterDistributeCache cache.DistributedCacheType[*do.Greeter], userDistributeCache cache.DistributedCacheType[*do.User], greeterLocalCache cache.LocalCache[string, *do.Greeter], userLocalCache cache.LocalCache[string, *do.User], rmqSender mq.MessageSender) *GreeterUsecase {
	return &GreeterUsecase{
		repo:                   repo,
		greeterDistributeCache: greeterDistributeCache,
		userDistributeCache:    userDistributeCache,
		greeterLocalCache:      greeterLocalCache,
		userLocalCache:         userLocalCache,
		rmqSender:              rmqSender,
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
	//present, err := uc.greeterLocalCache.GetIfPresent("test_local")
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(present)
	listener := NewDemoListener()
	message := &model.TopicMessage{
		Destination: "txTopic",
	}
	res, err := uc.rmqSender.SendMessageInTransaction(listener, message)
	if err != nil {
		return nil, err
	}
	fmt.Printf("事务发送结果: %v\n", res)
	return uc.repo.ListAll(ctx)
}

type DemoListener struct {
}

func NewDemoListener() *DemoListener {
	return &DemoListener{}
}

func (d *DemoListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	fmt.Printf("执行本地事务,%v", message)
	// 模拟执行本地事务
	return primitive.UnknowState
}

func (d *DemoListener) CheckLocalTransaction(ext *primitive.MessageExt) primitive.LocalTransactionState {
	// 当rocketmq无法确定本地事务状态时（收到的状态为非rollback或者commit时，或者超时未收到ack），会回调用客户端此方法
	fmt.Printf("检查本地事务,%v", ext)
	// 模拟检查本地事务状态
	return primitive.CommitMessageState
}
