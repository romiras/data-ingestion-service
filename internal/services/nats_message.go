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

// AckNATSMessage wraps a nats.Msg for acknowledgable pub/sub use.
type AckNATSMessage struct {
	NATSMessage
}

func (m *AckNATSMessage) Ack() error {
	// AckSync ensures the server has processed the ACK.
	return m.msg.AckSync()
}

func (m *AckNATSMessage) Nack() error {
	// Nak() tells JetStream to not redeliver this message.
	// If the stream is configured with a DLQ, it will go there.
	return m.msg.Nak()
}

func NewAckNATSMessage(natsMsg *nats.Msg) interfaces.AckMessage {
	return &AckNATSMessage{
		NATSMessage: NATSMessage{msg: natsMsg},
	}
}
