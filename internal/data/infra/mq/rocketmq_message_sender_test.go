package mq

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	. "github.com/smartystreets/goconvey/convey"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/model"
	"im-server/internal/pkg/infra/mq"
	plog "im-server/pkg/log"
	"testing"
	"time"
)

type DemoListener struct {
}

func NewDemoListener() *DemoListener {
	return &DemoListener{}
}

func (d *DemoListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	fmt.Printf("执行本地事务,%v", message)
	// 模拟执行本地事务
	time.Sleep(1 * time.Minute)
	return primitive.UnknowState
}

func (d *DemoListener) CheckLocalTransaction(ext *primitive.MessageExt) primitive.LocalTransactionState {
	// 当rocketmq无法确定本地事务状态时（收到的状态为非rollback或者commit时，或者超时未收到ack），会回调用客户端此方法
	fmt.Printf("检查本地事务,%v", ext)
	// 模拟检查本地事务状态
	return primitive.CommitMessageState
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
	sender, err := NewRocketMQMessageSender(config)
	if err != nil {
		return nil, err
	}
	plog.NewLogger("test", "", 100, 10, 10, plog.WithLevel("debug"))
	return sender, nil
}

func TestSyncSend(t *testing.T) {
	Convey("测试同步发送", t, func() {
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

func TestSendInTx(t *testing.T) {
	Convey("测试事务消息发送", t, func() {
		sender, err := InitRocketMQSender()
		So(err, ShouldBeNil)
		So(sender, ShouldNotBeNil)
		message := &model.TopicMessage{
			Destination: "txTopic",
		}
		listener := NewDemoListener()
		res, err := sender.SendMessageInTransaction(listener, message)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
