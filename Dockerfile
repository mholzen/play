# Build stage
FROM golang:1.23-alpine AS builder

# Install git, build dependencies, and ALSA development libraries
RUN apk add --no-cache git build-base alsa-lib-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies and create dialout group
RUN apk add --no-cache ca-certificates alsa-lib libstdc++ && \
    addgroup -S dialout && \
    adduser -S appuser -G dialout

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Set ownership and permissions
RUN chown appuser:dialout /app/server && \
    chmod +x /app/server

# Switch to non-root user
USER appuser

# Copy any necessary static files or assets if needed
# COPY static/ static/

# Expose the port the server runs on
EXPOSE 8080

# Run the server
CMD ["./server"] 
