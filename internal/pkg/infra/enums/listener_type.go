package enums

import "fmt"

type IMListenerType int

// 2. 使用iota自动生成递增枚举值，显式定义Unknown为-1
const (
	UnknownListener        IMListenerType = -1 // 显式定义为-1，区别于正常枚举值
	AllListener            IMListenerType = 0  // 0
	PrivateMessageListener IMListenerType = 1  // 1
	GroupMessageListener   IMListenerType = 2  // 2
)

// 3. 枚举值与描述的映射表（包含Unknown）
var listenerTypeDesc = map[IMListenerType]string{
	UnknownListener:        "未知监听类型",
	AllListener:            "全部消息",
	PrivateMessageListener: "私聊消息",
	GroupMessageListener:   "群聊消息",
}

func (t IMListenerType) String() string {
	if desc, ok := listenerTypeDesc[t]; ok {
		return desc
	}
	return fmt.Sprintf("未知监听类型(%d)", t) // 处理未定义的非法值
}

func (t IMListenerType) Code() int {
	return int(t)
}
