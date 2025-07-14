package model

type TopicMessage struct {
	Destination string `json:"destination"`
}

func NewTopicMessage(destination string) *TopicMessage {
	return &TopicMessage{
		Destination: destination,
	}
}
