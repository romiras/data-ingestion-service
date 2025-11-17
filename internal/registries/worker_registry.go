package registries

import (
	"log"

	"github.com/spf13/viper"
	"play.ground/generic-data-collector/initializers"
	"play.ground/generic-data-collector/internal/interfaces"
	"play.ground/generic-data-collector/internal/services"
)

type WorkerAppRegistry struct {
	Config         *viper.Viper
	Consumer       interfaces.Consumer
	BatchProcessor *services.BatchProcessor
}

func NewWorkerAppRegistry() (*WorkerAppRegistry, error) {
	env := getEnv()
	config := initializers.NewConfig(env)

	natsUrl := config.GetString("NATS_URL")
	if natsUrl == "" {
		log.Printf("NATS_URL not set in config, using default: %s", DefaultNATSUrl)
		natsUrl = DefaultNATSUrl
	}

	natsConsumer, err := services.NewNATSConsumer(natsUrl)
	if err != nil {
		log.Fatalf("Failed to create NATS consumer: %v", err)
		return nil, err
	}

	batchProcessor := services.NewBatchProcessor(natsConsumer)

	return &WorkerAppRegistry{
		Config:         config,
		Consumer:       natsConsumer,
		BatchProcessor: batchProcessor,
	}, nil
}

// NewMockWorkerAppRegistry creates a WorkerAppRegistry with a MockConsumer for testing.
func NewMockWorkerAppRegistry() *WorkerAppRegistry {
	mockConsumer := services.NewMockConsumer()
	return &WorkerAppRegistry{
		Consumer:       mockConsumer,
		BatchProcessor: services.NewBatchProcessor(mockConsumer),
	}
}
