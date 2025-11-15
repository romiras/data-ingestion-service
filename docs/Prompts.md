# Prompts

## Initial prompt

We want to build a generic recommended directory structure holding basic implementations of application services in Go.
REST API server can be implemented in Gin-gonic.

Example project structure for reference:

* `cmd/server/main.go` - a main function of API server
* `cmd/consumer/main.go` - a main function of consumer
* `internal/` - contains internal implementation which is not shared with external Go packages.
* `internal/interfaces/` - contains definitions of interfaces
* `internal/interfaces/kv.go` - contains definitions of interface for K/V storage
* `internal/interfaces/producer.go` - contains definitions of interface for message producer
* `internal/interfaces/consumer.go` - contains definitions of interface for message consumer
* `internal/handlers/` - contains implementation of handlers layers, like for Gin-gonic.
* `internal/services/` - contains implementation of services layers
* `internal/services/mock_kv.go` - contains implementation of mock K/V service
* `internal/services/mock_producer.go` - contains implementation of mock producer
* `internal/services/mock_consumer.go` - contains implementation of mock consumer
* `internal/registries/` - contains application registries. `AppRegistry` is a structure holding interfaces and, optionally, interface of logger and telemetry. New registry initializes concrete implementations of interfaces.
* `pkg/` - contains packages shared with external Go packages.

Main functions initialize a new registry and can add goroutines for handlers of events.
**IMPORTANT**: All initializations of concrete implementations must be done inside `internal`! There should not be any initializers like SQL DB, K/V, MQ etc.
Main functions should stay slim and readable.

We want to create a directory structure for REST API server in Gin-gonic that accepts dummy JSON body with stub value, publishes it via producer to MQ broker asynchronously. Right after that returns HTTP status: "accepted" with JSON field status: "ok.
In addition, we want to create another program for consuming messages with a worker that uses implementation of consumer of MQ broker. A consumer's handler should just print message body to stdout.

## Prompt refinements

* Server should have `ServerAppRegistry` (w/ `Producer` only) and Consumer should have `WorkerAppRegistry` (w/ `Consumer` only). Both should use NATS. Show me a plan of changes.
* add API server and Consumer to @docker-compose.yml . use alpine
* is `go run ...` relevant for Docker Compose or we should create Dockerfile for them?
* let's use `/app` directory for all pre-built binaries in Docker
* generate `.gitignore`
