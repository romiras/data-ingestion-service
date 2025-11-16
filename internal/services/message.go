package services

import (
	"play.ground/generic-data-collector/internal/interfaces"
)

// a non-acknowledgable pub/sub message implementation
type NonAckPubSubMessage struct {
	data []byte
}

func (msg NonAckPubSubMessage) Data() []byte {
	return msg.data
}

func NewNonAckPubSubMessage(data []byte) interfaces.Message {
	return NonAckPubSubMessage{data: data}
}
