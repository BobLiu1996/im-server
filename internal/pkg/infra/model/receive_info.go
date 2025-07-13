package model

type IMReceiveInfo struct {
	*TopicMessage               // 嵌入TopicMessage实现继承[6](@ref)
	Cmd           int           `json:"cmd"`
	Sender        *IMUserInfo   `json:"sender"`
	Receivers     []*IMUserInfo `json:"receivers"`
	SendResult    bool          `json:"sendResult"`
	Data          interface{}   `json:"data"`
}

func NewReceiveInfo(cmd int, sender *IMUserInfo, data interface{}) *IMReceiveInfo {
	return &IMReceiveInfo{
		TopicMessage: &TopicMessage{},
		Cmd:          cmd,
		Sender:       sender,
		Receivers:    make([]*IMUserInfo, 0),
		SendResult:   true,
		Data:         data,
	}
}

func (r *IMReceiveInfo) GetDestination() string {
	if r.TopicMessage == nil {
		return ""
	}
	return r.TopicMessage.Destination
}
