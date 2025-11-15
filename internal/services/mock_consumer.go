package services

import (
	"log"
	"time"
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
func (m *MockConsumer) Subscribe(topic string) (<-chan []byte, error) {
	log.Printf("MOCK CONSUMER: Subscribing to topic '%s'\n", topic)
	// Simulate receiving messages
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				msg := []byte(`{"data": "hello from mock consumer"}`)
				m.messages <- msg
			case <-m.done:
				return
			}
		}
	}()
	return m.messages, nil
}

// Close simulates closing the consumer.
func (m *MockConsumer) Close() error {
	log.Println("MOCK CONSUMER: Closed.")
	close(m.done)
	close(m.messages)
	return nil
}
