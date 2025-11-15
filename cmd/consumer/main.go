package main

import (
	"fmt"
	"play.ground/generic-data-collector/internal/registries"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	registry, err := registries.NewRegistry()
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}
	defer registry.Close()

	msgChan, err := registry.Consumer.Subscribe("metrics")
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Consumer worker started. Waiting for messages...")

	go func() {
		for msg := range msgChan {
			fmt.Printf("Received message: %s\n", string(msg))
		}
		log.Println("Message channel closed. Exiting.")
		quit <- syscall.SIGTERM // Trigger shutdown if channel closes
	}()

	<-quit
	log.Println("Shutting down consumer worker...")
}

