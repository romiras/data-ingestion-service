# Data ingestion service

This project is an example of a generic data ingestion service built in Go, designed for asynchronous, containerized data processing. It features a REST API server for receiving metrics and a worker consumer for processing them, all connected via a NATS message queue.

The architecture is designed for extensibility, utilizing interfaces for the producer and consumer components and a central registry for dependency management.

## Architecture Overview

| Component | Role | Technology |
| :--- | :--- | :--- |
| **Server** | REST API Producer | Gin-gonic, NATS |
| **Consumer** | Worker/Processor | NATS |
| **Messaging** | Message Queue | NATS |

## Prerequisites

* [Go](https://golang.org/dl/) (version 1.20 or higher)
* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Building and Running

The project uses Go modules for dependency management and Docker Compose to orchestrate the required services (NATS, Server, Consumer).

### 1. Run Infrastructure Services (NATS)

Start the NATS message queue and any other required infrastructure services using Docker Compose:

```bash
docker compose up -d nats
```

### 2. Run the Server (Producer)

The server exposes the REST API endpoint for data collection.

```bash
go run cmd/server/main.go
```

The server will start on port 8080.

### 3. Run the Consumer (Worker)

The consumer subscribes to the message queue and processes the received data.

```bash
go run cmd/consumer/main.go
```

The consumer will start and begin listening for messages on the "metrics" topic.

## Development

The project follows the standard Go project layout:

* `cmd`: Contains the main application entry points (`server` and `consumer`).
* `internal`: Contains the core business logic, including interfaces, handlers, registries, and services.

### Registries

The application uses separate registries for the server and worker to manage dependencies:

* `internal/registries/server_registry.go`: Manages the `Producer` dependency.
* `internal/registries/worker_registry.go`: Manages the `Consumer` dependency.

### Containerization

The project uses a multi-stage `Dockerfile` to create minimal, production-ready images for the server and consumer binaries.

To build and run all services using Docker Compose:

```bash
docker compose up --build
```

### Sending requests

#### Sending example data

```shell
curl --location 'http://localhost:8080/api/v1/metrics' \
--header 'Content-Type: application/json' \
--data '{
    "data": {
        "value": 123.45
    }
}'
```

Sample output:
```json
{
    "status": "ok"
}
```

## Next Steps

* Implement NATS JetStream for persistent and durable messaging (optionally).
* Add database integration for data storage.
