package mq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"im-server/internal/pkg/infra/model"
)

type MessageSender interface {
	Send(ctx context.Context, message *model.TopicMessage) (bool, error)
	SendMessageInTransaction(ctx context.Context, message *model.TopicMessage) (*primitive.TransactionSendResult, error)
}
