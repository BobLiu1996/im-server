package distribute

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"im-server/internal/conf"
	"im-server/internal/data"
	"im-server/internal/pkg/infra/cache"
	plog "im-server/pkg/log"
	"testing"
	"time"
)

const (
	testKeyPrefix = "test_distribute_cache"
)

type (
	User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	ID struct {
		Name    string
		Version int
	}
)

func InitRedisDistributeCacheService[T any]() (cache.DistributedCacheType[T], func(), error) {
	config := &conf.Data{
		Mysql: &conf.Data_MySql{
			Driver: "mysql",
			//Source: "root:root@tcp(192.168.5.134:3306)/test?charset=utf8mb4&parseTime=true&loc=Local",
			Source: "root:mystic@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Local",
		},
		Redis: &conf.Data_Redis{
			//Addr: "192.168.5.134:6379",
			Addr:      "localhost:6379",
			Db:        0,
			Pool:      250,
			IsCluster: false,
		},
		Debug: true,
	}
	d, f, err := data.NewData(config)
	if err != nil {
		return nil, nil, err
	}
	plog.NewLogger("test", "", 100, 10, 10, plog.WithLevel("debug"))
	redisCache := NewRedisDistributeCacheType[T](d)
	return redisCache, f, nil
}
func TestGetKey(t *testing.T) {
	Convey("测试getKey函数", t, func() {
		res := getKey(testKeyPrefix, &ID{
			Name:    "testId",
			Version: 1,
		})
		So(len(res), ShouldBeGreaterThan, 0)
	})
}

func TestGetResult(t *testing.T) {
	Convey("获取结果对象", t, func() {
		userStr := `{"name":"Alice","age":25}`
		res, err := GetResult[*User](userStr)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetResultList(t *testing.T) {
	Convey("获取结果列表", t, func() {
		userListStr := `[{"name":"Alice","age":25},{"name":"Bob","age":30}]`
		res, err := GetResultList[*User](userListStr)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 2)
	})
}

func TestSetCacheValue(t *testing.T) {
	Convey("设置缓存对象", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService[*User]()
		defer cleanup()
		So(err, ShouldBeNil)
		key := getKey(testKeyPrefix, "user1")
		value := &User{Name: "Bob", Age: 18}
		err = redisCache.Set(context.Background(), key, value)
		So(err, ShouldBeNil)
	})
}

func TestQueryWithPassThrough(t *testing.T) {
	Convey("以缓存穿透模式查询缓存对象-数据库中无数据", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService[*User]()
		defer cleanup()
		So(err, ShouldBeNil)
		id := "user2"
		// 模拟从数据库中拿到了空数据
		emptyFn := func(ctx context.Context, key any) (*User, error) {
			return nil, nil
		}
		d, err := redisCache.QueryWithPassThrough(context.Background(), testKeyPrefix, id, emptyFn, 5*time.Second)
		So(err, ShouldBeNil)
		res, _ := GetResult[*User](d)
		So(res, ShouldBeNil)
	})

	Convey("以缓存穿透模式查询缓存对象-数据库中存在数据", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService[*User]()
		defer cleanup()
		So(err, ShouldBeNil)
		id := "user3"
		// 模拟从数据库获取到非空数据
		noEmptyFn := func(ctx context.Context, key any) (*User, error) {
			user := &User{Name: "Alice", Age: 25}
			return user, nil
		}
		d, err := redisCache.QueryWithPassThrough(context.Background(), testKeyPrefix, id, noEmptyFn, 20*time.Second)
		So(err, ShouldBeNil)
		res, _ := GetResult[*User](d)
		So(res, ShouldNotBeNil)
	})
}

func TestQueryWithPassThroughList(t *testing.T) {
	Convey("以缓存穿透模式查询缓存对象列表-数据库中无数据", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService[*User]()
		defer cleanup()
		So(err, ShouldBeNil)
		id := "user4"
		// 模拟从数据库中拿到了空数据
		emptyFn := func(ctx context.Context, key any) ([]*User, error) {
			return nil, nil
		}
		d, err := redisCache.QueryWithPassThroughList(context.Background(), testKeyPrefix, id, emptyFn, 5*time.Second)
		So(err, ShouldBeNil)
		res, _ := GetResultList[*User](d)
		So(res, ShouldBeNil)
	})

	Convey("以缓存穿透模式查询缓存对象列表-数据库中存在数据", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService[*User]()
		defer cleanup()
		So(err, ShouldBeNil)
		id := "user5"
		// 模拟从数据库获取到非空数据
		noEmptyFn := func(ctx context.Context, key any) ([]*User, error) {
			userList := []*User{
				{Name: "Bob", Age: 30},
				{Name: "Charlie", Age: 35},
			}
			return userList, nil
		}
		d, err := redisCache.QueryWithPassThroughList(context.Background(), testKeyPrefix, id, noEmptyFn, 20*time.Second)
		So(err, ShouldBeNil)
		res, _ := GetResultList[*User](d)
		So(res, ShouldNotBeNil)
	})
}
