# Multi-stage build for production
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY ./server/go.mod ./server/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY ./server ./

# Copy migrations
COPY ./migrations ./migrations

RUN rm -rf .env .env.example

# Build the binary
RUN go build -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8000

# Run the application
CMD ["./main"]
