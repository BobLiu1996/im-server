package constants

const (
	// UserIdSplit *********************分隔符相关********************
	UserIdSplit     = "," // 用户ID分隔符
	RedisKeySplit   = ":" // Redis Key的分隔符
	MessageKeySplit = "_" // 发送消息的Key

	// OnlineTimeoutSeconds *******************基本信息相关*******************
	OnlineTimeoutSeconds = 600 // 在线状态过期时间，默认10分钟
	AllowRecallSecond    = 300 // 消息允许撤回时间，默认5分钟

	// ImMaxServerId *******************Redis相关*******************
	ImMaxServerId         = "im:max_server_id"           // bh-im-server最大ID
	ImUserServerId        = "im:user:server_id"          // 用户连接的IM-server ID
	ImGroupReadedPosition = "im:readed:group:position"   // 已读群聊消息位置
	ImWebrtcSession       = "im:webrtc:session"          // webrtc会话信息
	ImCache               = "im:cache:"                  // 缓存前缀
	ImCacheFriend         = ImCache + "friend"           // 是否好友缓存
	ImCacheGroup          = ImCache + "group"            // 群聊信息缓存
	ImCacheGroupMemberId  = ImCache + "group_member_ids" // 群聊成员ID缓存

	// MsgKey *******************RocketMQ相关*******************
	MsgKey                        = "message"                           // 消息key
	ImMessagePrivateQueue         = "im_message_private"                // 未读私聊队列
	ImMessagePrivateNullQueue     = "im_null_private"                   // 未读私聊空队列
	ImMessagePrivateConsumerGroup = "im_message_private_consumer_group" // 私聊消费分组
	ImMessageGroupQueue           = "im_message_group"                  // 未读群聊队列
	ImMessageGroupNullQueue       = "im_null_group"                     // 未读群聊空队列
	ImMessageGroupConsumerGroup   = "im_message_group_consumer_group"   // 群聊消费分组
	ImResultPrivateQueue          = "im_result_private"                 // 私聊结果队列
	ImResultPrivateConsumerGroup  = "im_result_private_consumer_group"  // 私聊结果分组
	ImResultGroupQueue            = "im_result_group"                   // 群聊结果队列
	ImResultGroupConsumerGroup    = "im_result_group_consumer_group"    // 群聊结果分组

	// UserId *******************Channel连接相关*******************
	UserId           = "USER_ID"         // 用户ID
	TerminalType     = "TERMINAL_TYPE"   // 终端类型
	HeartbeatTimes   = "HEARTBEAT_TIMES" // 心跳次数
	MinReadableBytes = 4                 // 最小读取字节数

	// MaxImageSize *******************平台相关*************************
	MaxImageSize   = 5 * 1024 * 1024  // 最大图片上传大小(5MB)
	MaxFileSize    = 10 * 1024 * 1024 // 最大文件上传大小(10MB)
	MaxGroupMember = 500              // 群聊最大人数
)

const (
	DistributedCacheRedisServiceKey = "distributed_cache_redis_service" // Redis服务key
	ImServerGroupBeanName           = "IMServerGroup"                   // IM服务组Bean名称
)
