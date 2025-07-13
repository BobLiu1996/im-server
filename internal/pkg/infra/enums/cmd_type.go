package enums

type IMCmdType int

// 2. 枚举常量声明（使用 iota 自动生成递增常量值）
const (
	ErrorCmd          IMCmdType = -1
	LoginCmd          IMCmdType = 0
	HeartBeatCmd      IMCmdType = 1
	ForceLogoutCmd    IMCmdType = 2
	PrivateMessageCmd IMCmdType = 3
	GroupMessageCmd   IMCmdType = 4
)

// 3. 枚举值到描述的映射表
var cmdTypeDesc = map[IMCmdType]string{
	LoginCmd:          "登录",
	HeartBeatCmd:      "心跳",
	ForceLogoutCmd:    "强制下线",
	PrivateMessageCmd: "私聊消息",
	GroupMessageCmd:   "群发消息",
}

func (c IMCmdType) String() string {
	if desc, ok := cmdTypeDesc[c]; ok {
		return desc
	}
	return "未知命令"
}

func (c IMCmdType) Code() int {
	return int(c)
}

func FromCode(code int) IMCmdType {
	if code >= 0 && code <= len(cmdTypeDesc) {
		return IMCmdType(code)
	}
	// 实际项目中建议返回error，此处简化为返回-1
	return ErrorCmd
}
