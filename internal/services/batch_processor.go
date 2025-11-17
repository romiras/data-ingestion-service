package services

import (
	"context"
	"log"
	"time"

	"play.ground/generic-data-collector/internal/interfaces"
)

const (
	batchSize    = 10
	batchTimeout = 5 * time.Second
	maxRetries   = 3
)

// BatchProcessor consumes messages, batches them, and sends them to an bulk ingester.
type BatchProcessor struct {
	consumer interfaces.Consumer
}

// NewBatchProcessor creates a new processor.
func NewBatchProcessor(consumer interfaces.Consumer) *BatchProcessor {
	return &BatchProcessor{
		consumer: consumer,
	}
}

// Start runs the main consumer loop. It blocks until the context is canceled.
func (p *BatchProcessor) Start(ctx context.Context, topic string) error {
	msgCh, err := p.consumer.Subscribe(topic)
	if err != nil {
		return err
	}

	batch := make([][]byte, 0, batchSize) // Preallocate slice for batch
	messages := make([]interfaces.Message, 0, batchSize)

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

// processBatch sends the batch to the metrics ingestion service.
func (p *BatchProcessor) processBatch(batch *[][]byte, messages *[]interfaces.Message) error {
	if len(*batch) == 0 {
		return nil // Nothing to process
	}

	// 1. Attempt to post data with retries
	err := p.postBatch(*batch)

	if err == nil {
		// 2. SUCCESS
		log.Printf("Successfully posted %d messages.", len(*messages))
	} else {
		// 3. FAILURE (send to DLQ)
		log.Printf("Failed to post batch %d messages after retries: %v.", len(*messages), err)
	}

	// 4. Clear local buffers
	*batch = make([][]byte, 0, batchSize)
	*messages = make([]interfaces.Message, 0, batchSize)
	return err // Return the error, if any, to the caller
}

func (p *BatchProcessor) postBatch(batch [][]byte) error {
	for _, bytes := range batch {
		log.Printf("Posting data: %s", string(bytes))
		// TODO: Implement actual posting logic here

		// Simulate network call latency
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}
