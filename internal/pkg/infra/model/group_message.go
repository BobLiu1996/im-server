package model

import "im-server/internal/pkg/infra/enums"

type IMGroupMessage[T any] struct {
	Sender           *IMUserInfo `json:"sender"`
	ReceiveIDs       []int64     `json:"receiveIds"`
	ReceiveTerminals []int       `json:"receiveTerminals"`
	SendToSelf       bool        `json:"sendToSelf"`
	SendResult       bool        `json:"sendResult"`
	Data             T           `json:"data"`
}

func NewGroupMessage[T any](
	sender *IMUserInfo,
	data T,
) *IMGroupMessage[T] {
	return &IMGroupMessage[T]{
		Sender:           sender,
		ReceiveIDs:       make([]int64, 0),            // 初始化空切片而非nil
		ReceiveTerminals: enums.IMTerminalTypeCodes(), // 假设已实现枚举方法
		SendToSelf:       true,                        // 显式设置默认值
		SendResult:       true,
		Data:             data,
	}
}
