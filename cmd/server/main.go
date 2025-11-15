package main

import (
	"play.ground/generic-data-collector/internal/registries"
	"play.ground/generic-data-collector/internal/routes"
	"log"
)

func main() {
	registry, err := registries.NewRegistry()
	if err != nil {
		log.Fatalf("Failed to initialize registry: %v", err)
	}
	defer registry.Close()

	routes.Run(registry)
}
