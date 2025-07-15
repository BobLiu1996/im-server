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
	producer rocketmq.Producer
	mqCfg    *conf.Data_RocketMQ
}

func NewRocketMQMessageSender(dataCfg *conf.Data) (*RocketMQMessageSender, error) {
	var p rocketmq.Producer
	mqCfg := dataCfg.GetRocketMQ()
	if mqCfg != nil {
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

	}
	return &RocketMQMessageSender{
		producer: p,
		mqCfg:    mqCfg,
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

func (r *RocketMQMessageSender) SendMessageInTransaction(txListener primitive.TransactionListener, message *model.TopicMessage) (*primitive.TransactionSendResult, error) {
	txp, err := rocketmq.NewTransactionProducer(
		txListener,
		producer.WithNameServer(r.mqCfg.GetNameServers()),
		producer.WithGroupName(r.mqCfg.GetProducer().GetGroupName()),
	)
	if err != nil {
		return nil, err
	}
	if err := txp.Start(); err != nil {
		return nil, err
	}
	defer txp.Shutdown()
	msg, err := r.buildMessage(message)
	if err != nil {
		return nil, err
	}
	return txp.SendMessageInTransaction(context.Background(), msg)
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
