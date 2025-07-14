package mq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"im-server/internal/pkg/infra/model"
)

type MessageSender interface {
	Send(message *model.TopicMessage) (bool, error)
	SendMessageInTransaction(message *model.TopicMessage) (*primitive.TransactionSendResult, error)
}
