package services

import (
	"log"
	"sync"
	"time"

	"play.ground/generic-data-collector/internal/interfaces"

	"github.com/nats-io/nats.go"
)

// NATSJetStreamConsumer implements the Consumer interface using NATS JetStream.
type NATSJetStreamConsumer struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	subs []*nats.Subscription
	mu   sync.Mutex
}

// NewNATSJetStreamConsumer creates a new consumer that connects to NATS
// and gets a JetStream context.
func NewNATSJetStreamConsumer(url string, streamName string) (interfaces.AckConsumer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, err
	}

	streamInfo, _ := js.StreamInfo(streamName)
	if streamInfo == nil {
		log.Printf("no stream found, creating %q stream", streamName)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"metrics"},
		})
		if err != nil {
			return nil, err
		}
	}

	return &NATSJetStreamConsumer{conn: nc, js: js}, nil
}

// Subscribe creates a durable, manual-ack JetStream subscription.
func (c *NATSJetStreamConsumer) Subscribe(topic string) (<-chan interfaces.AckMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Channel for nats.Msg from the NATS library
	natsMsgCh := make(chan *nats.Msg, 64)

	// Create a durable consumer. All pods in your K8s deployment
	// should use the SAME durable name to act as a competing consumer group.
	sub, err := c.js.ChanSubscribe(topic, natsMsgCh,
		nats.Durable("my-durable-consumer"), // <-- Key for K8s scaling
		nats.ManualAck(),                    // <-- CRITICAL for reliability
		nats.AckWait(time.Second*30),        // <-- How long to wait before redelivery
	)
	if err != nil {
		return nil, err
	}
	c.subs = append(c.subs, sub)

	// Channel for interfaces.AckMessage to return to the caller
	dataCh := make(chan interfaces.AckMessage, 64)

	// Goroutine to wrap nats.Msg into our interface
	go func() {
		defer close(dataCh)
		for msg := range natsMsgCh {
			dataCh <- NewAckNATSMessage(msg)
		}
	}()

	return dataCh, nil
}

// Close unsubscribes and closes the connection.
func (c *NATSJetStreamConsumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var firstErr error
	for _, sub := range c.subs {
		// Use Unsubscribe() for durable consumers to stop receiving,
		// but leave the consumer on the server.
		if err := sub.Unsubscribe(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	c.subs = nil
	c.conn.Close()
	return firstErr
}
