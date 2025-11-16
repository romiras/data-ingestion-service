package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"play.ground/generic-data-collector/internal/interfaces"
	"play.ground/generic-data-collector/internal/registries"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

// run contains all the application logic and returns an error if startup fails
func run() error {
	registry, err := registries.NewWorkerAppRegistry()
	if err != nil {
		return fmt.Errorf("failed to create registry: %w", err)
	}

	// Ensure consumer is closed on exit
	defer func() {
		if err := registry.Consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	// Create root context that will be cancelled on shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	// Start consuming
	if err := consumeMessages(ctx, registry.Consumer, "metrics"); err != nil {
		return fmt.Errorf("message consumption failed: %w", err)
	}

	log.Println("Graceful shutdown complete.")
	return nil
}

// consumeMessages starts consuming from the topic and returns when context is cancelled
func consumeMessages(ctx context.Context, consumer interfaces.Consumer, topic string) error {
	msgCh, err := consumer.Subscribe(topic)
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
