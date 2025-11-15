package services

import (
	"log"
)

// MockProducer is a mock implementation of the Producer interface.
type MockProducer struct{}

// NewMockProducer creates a new MockProducer.
func NewMockProducer() *MockProducer {
	return &MockProducer{}
}

// Publish simulates publishing a message by printing it to the console.
func (m *MockProducer) Publish(topic string, message []byte) error {
	log.Printf("MOCK PRODUCER: Publishing to topic '%s': %s\n", topic, string(message))
	return nil
}

// Close simulates closing the producer.
func (m *MockProducer) Close() error {
	log.Println("MOCK PRODUCER: Closed.")
	return nil
}
