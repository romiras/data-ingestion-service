package interfaces

// Consumer defines the interface for receiving messages from a pub/sub system.
type Consumer interface {
	// Subscribe starts listening to a given channel and returns a Go channel
	// from which messages can be read.
	Subscribe(topic string) (<-chan []byte, error)
	// Close stops the consumer and cleans up any underlying resources.
	Close() error
}
