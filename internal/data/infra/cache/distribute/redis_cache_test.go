package distribute

import (
	"context"
	"encoding/json"
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

func (u *User) MarshalBinary() ([]byte, error) {
	d, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func InitRedisDistributeCacheService() (cache.DistributedCache, func(), error) {
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
	redisCache := NewRedisDistributeCacheService(d)
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

func TestSetCacheValue(t *testing.T) {
	Convey("设置缓存对象", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService()
		defer cleanup()
		So(err, ShouldBeNil)
		key := getKey(testKeyPrefix, "user1")
		value := &User{Name: "Bob", Age: 18}
		err = redisCache.Set(context.Background(), key, value)
		So(err, ShouldBeNil)
	})
}

func TestQueryWithPassThrough(t *testing.T) {
	Convey("以缓存穿透模式查询缓存对象", t, func() {
		redisCache, cleanup, err := InitRedisDistributeCacheService()
		defer cleanup()
		So(err, ShouldBeNil)
		id := "user2"
		// 模拟从数据库获取到非空数据
		fn := func(ctx context.Context, key any) (any, error) {
			user := &User{Name: "Alice", Age: 25}
			return user, nil
		}
		// 模拟从数据库中拿到了空数据
		//fn := func(ctx context.Context, key any) (any, error) {
		//	return nil, nil
		//}
		d, err := redisCache.QueryWithPassThrough(context.Background(), testKeyPrefix, id, fn, 10*time.Second)
		So(err, ShouldBeNil)
		res := getResult[*User](d)
		So(res, ShouldNotBeNil)
	})
}
