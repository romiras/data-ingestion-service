## Project Overview

This project is a generic data collector built in Go. It consists of two main components: a REST API server and a consumer. The server exposes an endpoint for receiving metrics data, which it then publishes to a NATS message queue. The consumer subscribes to the message queue, receives the data, and prints it to the console.

The project is structured to be easily extensible, with interfaces for the producer and consumer, and a central registry for managing dependencies. It uses Gin-gonic for the web server and the official NATS Go client for message queuing.

## Building and Running

The project uses Go modules for dependency management and Docker Compose for running the required services.

### Running the Services

To run the NATS and Redis services, use the following command:

```bash
docker-compose up -d
```

### Running the Server

To run the API server, use the following command:

```bash
go run cmd/server/main.go
```

The server will start on port 8080.

### Running the Consumer

To run the consumer, use the following command:

```bash
go run cmd/consumer/main.go
```

The consumer will start and begin listening for messages on the "metrics" topic.

## Development Conventions

The project follows the standard Go project layout, with the main applications in the `cmd` directory and the internal packages in the `internal` directory.

### Interfaces

The project makes extensive use of interfaces to decouple the different components. The `internal/interfaces` directory contains the interfaces for the producer and consumer.

### Registry

The project uses a central registry to manage dependencies. The `internal/registries` directory contains the registry, which is responsible for creating and holding the producer and consumer services.

### Mock Implementations

The project provides mock implementations of the producer and consumer for testing purposes. These are located in the `internal/services` directory. To use the NATS implementations, you will need to update the `internal/registries/registry.go` file to use the `NewNATSProducer` and `NewNATSConsumer` functions.
