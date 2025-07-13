package enums

import "fmt"

type IMSendCode int

const (
	SendSuccess     IMSendCode = 0    // 发送成功
	NotOnline       IMSendCode = 1    // 对方当前不在线
	NotFoundChannel IMSendCode = 2    // 未找到对方的channel
	UnknownError    IMSendCode = 9999 // 未知异常
)

var sendCodeDesc = map[IMSendCode]string{
	SendSuccess:     "发送成功",
	NotOnline:       "对方当前不在线",
	NotFoundChannel: "未找到对方的channel",
	UnknownError:    "未知异常",
}

func (c IMSendCode) String() string {
	if desc, ok := sendCodeDesc[c]; ok {
		return desc
	}
	return fmt.Sprintf("未知发送状态(%d)", c) // 处理未定义的非法值
}

func (c IMSendCode) Code() int {
	return int(c)
}
