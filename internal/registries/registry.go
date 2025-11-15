package registries

import (
	"play.ground/generic-data-collector/internal/interfaces"
	"play.ground/generic-data-collector/internal/services"
)

// AppRegistry holds instances of all application-wide interfaces.
type AppRegistry struct {
	Producer interfaces.Producer
	Consumer interfaces.Consumer
}

// NewRegistry creates and returns an AppRegistry with mock implementations.
func NewRegistry() (*AppRegistry, error) {
	producer := services.NewMockProducer()
	consumer := services.NewMockConsumer()

	return &AppRegistry{
		Producer: producer,
		Consumer: consumer,
	}, nil
}

// Close cleans up resources used by the registry.
func (r *AppRegistry) Close() {
	if r.Producer != nil {
		r.Producer.Close()
	}
	if r.Consumer != nil {
		r.Consumer.Close()
	}
}
