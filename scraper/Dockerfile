# Stage 1: Build the Go application with librdkafka for ARM64
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application for ARM64 with CGO_ENABLED=1 to link against librdkafka
RUN go build -o kafka-producer .

# Stage 2: Run the Go application in a minimal container
FROM alpine:latest

# Copy the built Go binary from the builder stage
COPY --from=builder /app/kafka-producer /kafka-producer

# Command to run the executable
ENTRYPOINT ["/kafka-producer"]
