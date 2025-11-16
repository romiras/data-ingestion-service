package services

import (
	"github.com/nats-io/nats.go"
	"play.ground/generic-data-collector/internal/interfaces"
)

type (
	NATSMessage struct {
		msg *nats.Msg // Embed or reference the underlying NATS message
	}
)

func (m NATSMessage) Data() []byte {
	return m.msg.Data // Direct delegation; customize if you need to include Subject/Header
}

func NewNATSMessage(natsMsg *nats.Msg) interfaces.Message {
	// Optional: Create a copy if you don't want to mutate the original nats.Msg
	// cloned := *natsMsg // Shallow copy; deep copy Data if needed for safety
	return NATSMessage{msg: natsMsg}
}
