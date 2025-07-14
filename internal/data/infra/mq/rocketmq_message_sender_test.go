package mq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	. "github.com/smartystreets/goconvey/convey"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/model"
	"im-server/internal/pkg/infra/mq"
	plog "im-server/pkg/log"
	"testing"
)

type DemoListener struct {
}

func NewDemoListener() *DemoListener {
	return &DemoListener{}
}

func (d *DemoListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	//TODO implement me
	panic("implement me")
}

func (d *DemoListener) CheckLocalTransaction(ext *primitive.MessageExt) primitive.LocalTransactionState {
	//TODO implement me
	panic("implement me")
}

func InitRocketMQSender() (mq.MessageSender, error) {
	config := &conf.Data{
		Mysql: &conf.Data_MySql{
			Driver: "mysql",
			Source: "root:root@tcp(192.168.5.134:3306)/test?charset=utf8mb4&parseTime=true&loc=Local",
			//Source: "root:mystic@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Local",
		},
		Redis: &conf.Data_Redis{
			Addr: "192.168.5.134:6379",
			//Addr:      "localhost:6379",
			Db:        0,
			Pool:      250,
			IsCluster: false,
		},
		RocketMQ: &conf.Data_RocketMQ{
			NameServers: []string{"192.168.5.134:9876"},
			Producer: &conf.Data_RocketMQ_Producer{
				GroupName: "test",
			},
		},
	}
	sender, err := NewRocketMQMessageSender(config, NewDemoListener())
	if err != nil {
		return nil, err
	}
	plog.NewLogger("test", "", 100, 10, 10, plog.WithLevel("debug"))
	return sender, nil
}
func TestGetKey(t *testing.T) {
	Convey("测试getKey函数", t, func() {
		sender, err := InitRocketMQSender()
		So(err, ShouldBeNil)
		So(sender, ShouldNotBeNil)
		message := &model.TopicMessage{
			Destination: "input",
		}
		ok, err := sender.Send(message)
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)
	})
}
