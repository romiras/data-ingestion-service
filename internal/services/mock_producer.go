package services

import (
	"log"
	"sync"
)

// MockProducer is a mock implementation of the Producer interface.
type MockProducer struct {
	PublishedData    []byte
	PublishedChannel string
	mu               sync.Mutex
}

// NewMockProducer creates a new MockProducer.
func NewMockProducer() *MockProducer {
	return &MockProducer{}
}

// Publish simulates publishing a message by printing it to the console and storing the data.
func (m *MockProducer) Publish(topic string, message []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.PublishedChannel = topic
	m.PublishedData = message

	log.Printf("MOCK PRODUCER: Publishing to topic '%s': %s\n", topic, string(message))
	return nil
}

// Close simulates closing the producer.
func (m *MockProducer) Close() error {
	log.Println("MOCK PRODUCER: Closed.")
	return nil
}
