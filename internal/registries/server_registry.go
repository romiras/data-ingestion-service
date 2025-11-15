package registries

import (
	"log"

	"github.com/spf13/viper"
	"play.ground/generic-data-collector/initializers"
	"play.ground/generic-data-collector/internal/interfaces"
	"play.ground/generic-data-collector/internal/services"
)

const DefaultNATSUrl = "nats://localhost:4222"

type ServerAppRegistry struct {
	Config   *viper.Viper
	Producer interfaces.Producer
}

func NewServerAppRegistry() (*ServerAppRegistry, error) {
	env := getEnv()
	config := initializers.NewConfig(env)
	natsUrl := config.GetString("NATS_URL")
	if natsUrl == "" {
		log.Printf("NATS_URL not set in config, using default: %s", DefaultNATSUrl)
		natsUrl = DefaultNATSUrl
	}

	natsProducer, err := services.NewNATSProducer(natsUrl)
	if err != nil {
		log.Fatalf("Failed to create NATS producer: %v", err)
		return nil, err
	}

	return &ServerAppRegistry{
		Config:   config,
		Producer: natsProducer,
	}, nil
}
