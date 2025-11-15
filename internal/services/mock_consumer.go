package services

import (
	"log"
)

// MockConsumer is a mock implementation of the Consumer interface.
type MockConsumer struct {
	messages chan []byte
	done     chan struct{}
}

// NewMockConsumer creates a new MockConsumer.
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		messages: make(chan []byte, 10),
		done:     make(chan struct{}),
	}
}

// Subscribe simulates subscribing to a topic.
// It returns the internal message channel, allowing tests to manually push messages.
func (m *MockConsumer) Subscribe(topic string) (<-chan []byte, error) {
	log.Printf("MOCK CONSUMER: Subscribing to topic '%s'\n", topic)
	// No ticker, messages must be pushed manually in tests.
	return m.messages, nil
}

// SendMessage allows tests to manually inject a message into the consumer's channel.
func (m *MockConsumer) SendMessage(message []byte) {
	m.messages <- message
}

// Close simulates closing the consumer.
func (m *MockConsumer) Close() error {
	log.Println("MOCK CONSUMER: Closed.")
	close(m.done)
	close(m.messages)
	return nil
}
