package interfaces

type (
	Message interface {
		// Data returns the message payload.
		Data() []byte
	}

	// Consumer defines the interface for receiving messages from a pub/sub system.
	Consumer interface {
		// Subscribe starts listening to a given channel and returns a Go channel
		// from which messages can be read.
		Subscribe(topic string) (<-chan Message, error)
		// Close stops the consumer and cleans up any underlying resources.
		Close() error
	}
)
