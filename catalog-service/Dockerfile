# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Ent code
RUN go run entgo.io/ent/cmd/ent generate ./ent/schema

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o catalog-service ./cmd/catalog

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/catalog-service .

# Expose ports
EXPOSE 8080
EXPOSE 9090

# Command to run
CMD ["./catalog-service"]