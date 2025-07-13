package enums

import "fmt"

type IMTerminalType int

const (
	Web IMTerminalType = 0 // web
	App IMTerminalType = 1 // app
)

var terminalTypeDesc = map[IMTerminalType]string{
	Web: "web",
	App: "app",
}

func (t IMTerminalType) String() string {
	if desc, ok := terminalTypeDesc[t]; ok {
		return desc
	}
	return fmt.Sprintf("未知终端类型(%d)", t) // 处理未定义的非法值
}

func (t IMTerminalType) Code() int {
	return int(t)
}

func IMTerminalTypeCodes() []int {
	res := make([]int, 0)
	for terminalType := range terminalTypeDesc {
		res = append(res, int(terminalType))
	}
	return res
}
