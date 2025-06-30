# Use the official Golang image (stable version)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum if present
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go server
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Use a minimal image for running
FROM alpine:3.19


WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/ .

# Expose port 8080 (change if your server uses a different port)
EXPOSE 8080

# Run the server
ENTRYPOINT ["./main"]