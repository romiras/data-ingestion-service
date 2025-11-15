FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server and consumer applications
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/consumer cmd/consumer/main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Copy the compiled binaries from the builder stage
COPY --from=builder /app/server /app/server
COPY --from=builder /app/consumer /app/consumer
COPY config /app/config

WORKDIR /app
