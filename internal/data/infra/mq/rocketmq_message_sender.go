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
	dataCfg    *conf.Data
	producer   rocketmq.Producer
	txProducer rocketmq.TransactionProducer
}

func NewRocketMQMessageSender(dataConf *conf.Data, txListener primitive.TransactionListener) (*RocketMQMessageSender, error) {
	rmqSender := &RocketMQMessageSender{
		dataCfg: dataConf,
	}
	if err := rmqSender.initRocketMQTxProducer(txListener); err != nil {
		return nil, err
	}
	if err := rmqSender.initRocketMQProducer(); err != nil {
		return nil, err
	}
	return rmqSender, nil
}

func (r *RocketMQMessageSender) initRocketMQProducer() error {
	var p rocketmq.Producer
	mqCfg := r.dataCfg.GetRocketMQ()
	if mqCfg != nil {
		var err error
		p, err = rocketmq.NewProducer(
			producer.WithRetry(1),
			producer.WithNsResolver(primitive.NewPassthroughResolver(mqCfg.GetNameServers())),
			producer.WithGroupName(mqCfg.GetProducer().GetGroupName()),
		)
		if err != nil {
			return err
		}
		err = p.Start()
		if err != nil {
			return err
		}
		r.producer = p
	}
	return nil
}

func (r *RocketMQMessageSender) initRocketMQTxProducer(txListener primitive.TransactionListener) error {
	mqCfg := r.dataCfg.GetRocketMQ()
	txProducerCfg := mqCfg.GetTxProducer()
	if txProducerCfg != nil {
		var err error
		txp, err := rocketmq.NewTransactionProducer(
			txListener,
			producer.WithNsResolver(primitive.NewPassthroughResolver(mqCfg.GetNameServers())),
			producer.WithGroupName(txProducerCfg.GetGroupName()),
			producer.WithRetry(1),
		)
		if err != nil {
			return err
		}
		err = txp.Start()
		if err != nil {
			return err
		}
		r.txProducer = txp
	}
	return nil
}

func (r *RocketMQMessageSender) Send(ctx context.Context, message *model.TopicMessage) (bool, error) {
	err := r.producer.Start()
	if err != nil {
		return false, err
	}
	msg, err := r.buildMessage(message)
	if err != nil {
		return false, err
	}
	res, err := r.producer.SendSync(ctx, msg)
	if err != nil {
		return false, err
	}
	return res.Status == primitive.SendOK, nil
}

func (r *RocketMQMessageSender) SendMessageInTransaction(ctx context.Context, message *model.TopicMessage) (*primitive.TransactionSendResult, error) {
	msg, err := r.buildMessage(message)
	if err != nil {
		return nil, err
	}
	return r.txProducer.SendMessageInTransaction(ctx, msg)
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
