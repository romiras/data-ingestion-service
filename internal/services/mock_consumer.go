package services

import (
	"log"

	"play.ground/generic-data-collector/internal/interfaces"
)

// MockConsumer is a mock implementation of the Consumer interface.
type MockConsumer struct {
	messages chan interfaces.Message
	done     chan struct{}
}

// NewMockConsumer creates a new MockConsumer.
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		messages: make(chan interfaces.Message, 10),
		done:     make(chan struct{}),
	}
}

// Subscribe simulates subscribing to a topic.
// It returns the internal message channel, allowing tests to manually push messages.
func (m *MockConsumer) Subscribe(topic string) (<-chan interfaces.Message, error) {
	log.Printf("MOCK CONSUMER: Subscribing to topic '%s'\n", topic)
	return m.messages, nil
}

// SendMessage allows tests to manually inject a message into the consumer's channel.
func (m *MockConsumer) SendMessage(message interfaces.Message) {
	m.messages <- message
}

// Close simulates closing the consumer.
func (m *MockConsumer) Close() error {
	log.Println("MOCK CONSUMER: Closed.")
	close(m.done)
	close(m.messages)
	return nil
}
