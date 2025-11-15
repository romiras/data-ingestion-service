package services

import (
	"play.ground/generic-data-collector/internal/interfaces"

	"github.com/nats-io/nats.go"
)

// NATSProducer implements the Producer interface for sending messages to NATS.
type NATSProducer struct {
	conn *nats.Conn
}

// NewNATSProducer creates a new producer that connects to the given NATS URL.
func NewNATSProducer(url string) (interfaces.Producer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSProducer{conn: nc}, nil
}

// Publish sends a message to a specific topic in NATS.
func (p *NATSProducer) Publish(topic string, message []byte) error {
	return p.conn.Publish(topic, message)
}

// Close drains and closes the NATS connection.
func (p *NATSProducer) Close() error {
	return p.conn.Drain()
}
