package services

import (
	"log"

	"play.ground/generic-data-collector/internal/interfaces"
)

type MockAckConsumer struct {
	messages chan interfaces.AckMessage
	done     chan struct{}
}

func NewMockAckConsumer() *MockAckConsumer {
	return &MockAckConsumer{
		messages: make(chan interfaces.AckMessage, 10),
		done:     make(chan struct{}),
	}
}

func (m *MockAckConsumer) Subscribe(topic string) (<-chan interfaces.AckMessage, error) {
	log.Printf("MOCK ACK CONSUMER: Subscribing to topic '%s'\n", topic)
	return m.messages, nil
}

func (m *MockAckConsumer) SendMessage(message interfaces.AckMessage) {
	m.messages <- message
}

func (m *MockAckConsumer) Close() error {
	log.Println("MOCK ACK CONSUMER: Closed.")
	close(m.done)
	close(m.messages)
	return nil
}
