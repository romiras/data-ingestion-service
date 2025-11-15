package main

import (
	"log"

	"play.ground/generic-data-collector/internal/registries"
	"play.ground/generic-data-collector/internal/routes"
)

func main() {
	registry, err := registries.NewServerAppRegistry()
	if err != nil {
		log.Fatalf("Failed to initialize server registry: %v", err)
	}

	routes.Run(registry)
}
