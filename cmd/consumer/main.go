package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"play.ground/generic-data-collector/internal/registries"
)

const (
	TopicName = "metrics"
)

func main() {
	registry, err := registries.NewWorkerAppRegistry()
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}

	runner := func() error {
		if registry.Config.GetBool("RUN_WITH_BATCHES") {
			return runWithBatches(registry)
		}
		return run(registry)
	}

	if err := runner(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

// run contains all the application logic and returns an error if startup fails
func run(registry *registries.WorkerAppRegistry) error {
	// Ensure consumer is closed on exit
	defer func() {
		if err := registry.Consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	// Create root context that will be cancelled on shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, os.Interrupt)
	defer stop()

	log.Println("Consumer worker starting...")

	// Start consuming
	if err := consumeMessages(ctx, registry, TopicName); err != nil {
		return fmt.Errorf("message consumption failed: %w", err)
	}

	log.Println("Graceful shutdown complete.")
	return nil
}

// consumeMessages starts consuming from the topic and returns when context is cancelled
func consumeMessages(ctx context.Context, registry *registries.WorkerAppRegistry, topic string) error {
	msgCh, err := registry.Consumer.Subscribe(topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %q: %s", err, topic)
	}

	log.Printf("Consumer worker started. Subscribed to topic: %s", topic)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutdown signal received, stopping message consumption...")
				return
			case msg, ok := <-msgCh:
				if !ok {
					log.Println("Message channel closed by publisher.")
					return
				}
				// Process message (you can add error handling / retry logic here)
				fmt.Printf("Received message: %v\n", string(msg.Data()))
			}
		}
	}()

	// Wait for either context cancellation or goroutine exit
	wg.Wait()
	log.Println("Message consumer stopped.")
	return nil
}

func runWithBatches(registry *registries.WorkerAppRegistry) error {
	// Use context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start the batch processor in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := registry.BatchProcessor.Start(ctx, TopicName); err != nil {
			log.Printf("Batch processor exited with error: %v", err)
		} else {
			log.Println("Batch processor exited gracefully.")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	<-quit
	log.Println("Shutting down consumer worker...")

	// Signal the batch processor to stop
	cancel()

	// Wait for the processor to finish processing its final batch
	log.Println("Waiting for batch processor to shut down...")
	wg.Wait()

	// Now, safely close the consumer connection
	if err := registry.Consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}

	log.Println("Shutdown complete.")
	return nil
}
