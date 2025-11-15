package services

import (
	"sync"

	"play.ground/generic-data-collector/internal/interfaces"

	"github.com/nats-io/nats.go"
)

// NATSConsumer implements the Consumer interface for receiving messages from NATS.
type NATSConsumer struct {
	conn *nats.Conn
	subs []*nats.Subscription
	mu   sync.Mutex
}

// NewNATSConsumer creates a new consumer that connects to the given NATS URL.
func NewNATSConsumer(url string) (interfaces.Consumer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NATSConsumer{conn: nc}, nil
}

// Subscribe starts listening to a given topic and returns a Go channel
// from which message payloads can be read.
func (c *NATSConsumer) Subscribe(topic string) (<-chan []byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Channel for nats.Msg from the NATS library
	natsMsgCh := make(chan *nats.Msg, 64)

	sub, err := c.conn.ChanSubscribe(topic, natsMsgCh)
	if err != nil {
		return nil, err
	}
	c.subs = append(c.subs, sub)

	// Channel for []byte to return to the caller
	dataCh := make(chan []byte, 64)

	// Goroutine to transfer message data from nats.Msg channel to the data channel.
	go func() {
		defer close(dataCh)
		for msg := range natsMsgCh {
			dataCh <- msg.Data
		}
	}()

	return dataCh, nil
}

// Close unsubscribes from all topics and closes the NATS connection.
func (c *NATSConsumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var firstErr error
	for _, sub := range c.subs {
		if err := sub.Unsubscribe(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	c.subs = nil // Clear the subscriptions

	c.conn.Close() // nats.Conn.Close() does not return an error.

	return firstErr
}
