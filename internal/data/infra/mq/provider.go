package mq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/wire"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/mq"
)

var ProviderSet = wire.NewSet(ProvideRocketMQMessageSender, ProvideTxListener)

func ProvideRocketMQMessageSender(dataCfg *conf.Data, txListener primitive.TransactionListener) mq.MessageSender {
	sender, err := NewRocketMQMessageSender(dataCfg, txListener)
	if err != nil {
		return nil
	}
	return sender
}

func ProvideTxListener() primitive.TransactionListener {
	return NewDemoListener()
}
