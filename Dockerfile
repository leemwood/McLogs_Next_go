# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o mclogs-api ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and configs
COPY --from:builder /app/mclogs-api .
COPY --from:builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Run
CMD ["./mclogs-api"]
