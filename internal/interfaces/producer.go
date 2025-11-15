package interfaces

// Producer defines the interface for sending messages to a pub/sub system.
type Producer interface {
	// Publish sends a message to a specific channel/topic.
	Publish(topic string, message []byte) error
	// Close cleans up any underlying resources.
	Close() error
}
