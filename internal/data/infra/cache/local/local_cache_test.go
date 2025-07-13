package local

import (
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/types/known/durationpb"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/cache"
	plog "im-server/pkg/log"
	"testing"
)

var (
	testKey string = "test_local_cache"
)

type (
	User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)

func InitLocalCacheService[K comparable, V any]() (cache.LocalCache[K, V], error) {
	config := &conf.Data{
		LocalCache: &conf.Data_LocalCache{
			Expiration: &durationpb.Duration{
				Seconds: 10,
				Nanos:   0,
			},
			CleanupInterval: &durationpb.Duration{
				Seconds: 600,
				Nanos:   0,
			},
		},
	}
	// 创建分布式锁
	plog.NewLogger("test", "", 100, 10, 10, plog.WithLevel("debug"))
	localCache := NewLocalGoCache[K, V](config)
	return localCache, nil
}

func TestLocalCache(t *testing.T) {
	localCache, err := InitLocalCacheService[string, *User]()
	if err != nil {
		return
	}
	Convey("测试set本地缓存", t, func() {
		user := &User{
			Name: "bob",
			Age:  18,
		}
		localCache.Put(testKey, user)

		cache1, err := localCache.GetIfPresent(testKey)
		So(err, ShouldBeNil)
		So(cache1.Name, ShouldEqual, user.Name)
		So(cache1.Age, ShouldEqual, user.Age)

		localCache.Remove(testKey)

		cache2, err := localCache.GetIfPresent(testKey)
		So(err, ShouldBeNil)
		So(cache2, ShouldBeNil)

	})
}
