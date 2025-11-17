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

	// AckMessage represents a single consumable message,
	// abstracting the underlying broker's delivery mechanism.
	AckMessage interface {
		Message

		// Ack acknowledges the message, marking it as processed.
		Ack() error
		// Nack negatively acknowledges the message (e.g., send to DLQ).
		Nack() error
	}

	// Consumer defines the interface for receiving acknowlegable messages from a pub/sub system.
	AckConsumer interface {
		// Subscribe starts listening to a given topic and returns a Go channel
		// from which Message objects can be read.
		Subscribe(topic string) (<-chan AckMessage, error)
		// Close stops the consumer and cleans up any underlying resources.
		Close() error
	}
)
