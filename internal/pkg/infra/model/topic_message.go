package model

import (
	"encoding/json"
)

type TopicMessage struct {
	Destination string `json:"destination"`
}

func (m *TopicMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func NewTopicMessage(destination string) *TopicMessage {
	return &TopicMessage{
		Destination: destination,
	}
}
