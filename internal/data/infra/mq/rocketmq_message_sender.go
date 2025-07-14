package mq

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"im-server/internal/conf"
	"im-server/internal/pkg/infra/constants"
	"im-server/internal/pkg/infra/model"
)

type RocketMQMessageSender struct {
	producer   rocketmq.Producer
	txProducer rocketmq.TransactionProducer
}

func NewRocketMQMessageSender(dataCfg *conf.Data, txListener primitive.TransactionListener) (*RocketMQMessageSender, error) {
	var p rocketmq.Producer
	var txp rocketmq.TransactionProducer
	if mqCfg := dataCfg.GetRocketMQ(); mqCfg != nil {
		addr, err := primitive.NewNamesrvAddr(mqCfg.GetNameServers()...)
		if err != nil {
			return nil, err
		}
		p, err = rocketmq.NewProducer(
			producer.WithNameServer(addr),
			producer.WithRetry(2),
			producer.WithGroupName(mqCfg.GetProducer().GetGroupName()),
		)
		if err != nil {
			return nil, err
		}

		txp, err = rocketmq.NewTransactionProducer(
			txListener,
			producer.WithNameServer(addr),
			producer.WithGroupName(mqCfg.GetProducer().GetGroupName()),
		)
		if err != nil {
			return nil, err
		}
	}
	return &RocketMQMessageSender{
		producer:   p,
		txProducer: txp,
	}, nil
}

func (r *RocketMQMessageSender) Send(message *model.TopicMessage) (bool, error) {
	err := r.producer.Start()
	if err != nil {
		return false, err
	}
	defer r.producer.Shutdown()
	msg, err := r.buildMessage(message)
	if err != nil {
		return false, err
	}
	res, err := r.producer.SendSync(context.Background(), msg)
	if err != nil {
		return false, err
	}
	return res.Status == primitive.SendOK, nil
}

func (r *RocketMQMessageSender) SendMessageInTransaction(message *model.TopicMessage) (*primitive.TransactionSendResult, error) {
	r.txProducer.Start()
	defer r.txProducer.Shutdown()
	msg, err := r.buildMessage(message)
	if err != nil {
		return nil, err
	}
	return r.txProducer.SendMessageInTransaction(context.Background(), msg)
}

func (r *RocketMQMessageSender) buildMessage(msg *model.TopicMessage) (*primitive.Message, error) {
	payload := map[string]interface{}{
		constants.MsgKey: msg,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return primitive.NewMessage(msg.Destination, jsonData), nil
}
