package model

import "im-server/internal/pkg/infra/enums"

type IMPrivateMessage[T any] struct {
	Sender           *IMUserInfo `json:"sender"`
	ReceiveID        int64       `json:"receiveId"`
	ReceiveTerminals []int       `json:"receiveTerminals"`
	SendToSelf       bool        `json:"sendToSelf"`
	SendResult       bool        `json:"sendResult"`
	Data             T           `json:"data"`
}

func NewIMPrivateMessage[T any](
	sender *IMUserInfo,
	receiveId int64,
	data T,
) *IMPrivateMessage[T] {
	return &IMPrivateMessage[T]{
		Sender:           sender,
		ReceiveID:        receiveId,
		ReceiveTerminals: enums.IMTerminalTypeCodes(), // 假设已实现枚举的Codes()方法[1](@ref)
		SendToSelf:       true,                        // 显式设置默认值
		SendResult:       true,
		Data:             data,
	}
}
