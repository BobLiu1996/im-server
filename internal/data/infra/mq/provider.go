package mq

import (
	"github.com/google/wire"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/mq"
)

var ProviderSet = wire.NewSet(ProvideRocketMQMessageSender)

func ProvideRocketMQMessageSender(dataConfig *conf.Data) mq.MessageSender {
	sender, err := NewRocketMQMessageSender(dataConfig)
	if err != nil {
		return nil
	}
	return sender
}
