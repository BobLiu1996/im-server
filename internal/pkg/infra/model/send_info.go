package model

type IMSendInfo[T any] struct {
	Cmd  int `json:"cmd"`  // 命令类型（对应Java的Integer）
	Data T   `json:"data"` // 泛型数据
}

func NewIMSendInfo[T any](cmd int, data T) *IMSendInfo[T] {
	return &IMSendInfo[T]{
		Cmd:  cmd,
		Data: data,
	}
}
