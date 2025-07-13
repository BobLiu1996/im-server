package model

import (
	"encoding/json"
)

type IMSendResult[T any] struct {
	*TopicMessage
	Sender   *IMUserInfo `json:"sender"`
	Receiver *IMUserInfo `json:"receiver"`
	Code     int         `json:"code"`
	Data     T           `json:"data"`
}

func NewIMSendResult[T any](
	sender, receiver *IMUserInfo,
	code int,
	data T,
) *IMSendResult[T] {
	return &IMSendResult[T]{
		TopicMessage: &TopicMessage{},
		Sender:       sender,
		Receiver:     receiver,
		Code:         code,
		Data:         data,
	}
}

func (m *IMSendResult[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}
