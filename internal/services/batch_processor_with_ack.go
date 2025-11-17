package services

import (
	"context"
	"log"
	"time"

	"play.ground/generic-data-collector/internal/interfaces"
)

// BatchProcessor consumes messages, batches them, and sends them to an bulk ingester.
type BatchProcessorWithAck struct {
	consumer interfaces.AckConsumer
}

// NewBatchProcessor creates a new processor.
func NewBatchProcessorWithAck(consumer interfaces.AckConsumer) *BatchProcessorWithAck {
	return &BatchProcessorWithAck{
		consumer: consumer,
	}
}

// Start runs the main consumer loop. It blocks until the context is canceled.
func (p *BatchProcessorWithAck) Start(ctx context.Context, topic string) error {
	msgCh, err := p.consumer.Subscribe(topic)
	if err != nil {
		return err
	}

	batch := make([][]byte, 0, batchSize) // Preallocate slice for batch
	messages := make([]interfaces.AckMessage, 0, batchSize)

	ticker := time.NewTicker(batchTimeout)
	defer ticker.Stop()

	log.Println("Batch processor started. Waiting for messages...")

	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				log.Println("Message channel closed.")
				// Channel is closed - process final batch and exit
				return p.processBatch(&batch, &messages)
			}

			batch = append(batch, msg.Data())
			messages = append(messages, msg)

			if len(batch) >= batchSize {
				log.Printf("Flushing batch: size limit reached (%d)", len(batch))
				if err := p.processBatch(&batch, &messages); err != nil {
					log.Printf("ERROR processing batch: %v", err)
					// Note: processBatch messages NACKing, so we just log and continue.
				}
				ticker.Reset(batchTimeout)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				log.Printf("Flushing batch: timeout reached (%d)", len(batch))
				if err := p.processBatch(&batch, &messages); err != nil {
					log.Printf("ERROR processing batch: %v", err)
				}
			}

		case <-ctx.Done():
			log.Println("Shutdown signal received. Processing final batch...")
			// Context canceled, process final batch and exit
			return p.processBatch(&batch, &messages)
		}
	}
}

// processBatch sends the batch to the metrics ingestion service and messages ACK/NACK.
func (p *BatchProcessorWithAck) processBatch(batch *[][]byte, messages *[]interfaces.AckMessage) error {
	if len(*batch) == 0 {
		return nil // Nothing to process
	}

	// 1. Attempt to post data with retries
	err := p.postBatch(*batch)

	if err == nil {
		// 2. SUCCESS: ACK all messages
		log.Printf("Successfully posted %d messages. ACKing.", len(*messages))
		for _, msg := range *messages {
			if ackErr := msg.Ack(); ackErr != nil {
				log.Printf("WARNING: Failed to ACK message: %v", ackErr)
			}
		}
	} else {
		// 3. FAILURE: NACK all messages (send to DLQ)
		log.Printf("Failed to post batch after retries: %v. NACKing %d messages.", err, len(*messages))
		for _, msg := range *messages {
			if nackErr := msg.Nack(); nackErr != nil {
				log.Printf("WARNING: Failed to NACK message: %v", nackErr)
			}
		}
	}

	// 4. Clear local buffers
	*batch = make([][]byte, 0, batchSize)
	*messages = make([]interfaces.AckMessage, 0, batchSize)
	return err // Return the error, if any, to the caller
}

func (p *BatchProcessorWithAck) postBatch(batch [][]byte) error {
	for _, bytes := range batch {
		log.Printf("Posting data: %s", string(bytes))
		// Simulate network call latency
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}
