# Dockerfile
# Stage 1: Development dependencies
FROM golang:1.21.4 AS dev-deps
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN sed -i 's/go 1.23.1/go 1.21/' go.mod && \
    go mod download

# Stage 2: Builder
FROM golang:1.21.4 AS builder
WORKDIR /go/src/app
COPY --from=dev-deps /go/pkg /go/pkg
COPY . .
RUN sed -i 's/go 1.23.1/go 1.21/' go.mod && \
    CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" \
    -o apiserver ./cmd/apiserver/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies and debugging tools
RUN apt-get update && \
    apt-get install -y \
    ca-certificates \
    tzdata \
    wget \
    netcat-traditional \
    curl \
    default-mysql-client \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -r -s /bin/false appuser \
    && mkdir -p /app/logs \
    && chown -R appuser:appuser /app

# Copy binary and configs
COPY --from=builder /go/src/app/apiserver /app/
COPY --from=builder /go/src/app/config /app/config/

# Create entrypoint script with database connection check
RUN echo '#!/bin/sh' > /app/docker-entrypoint.sh && \
    echo 'set -e' >> /app/docker-entrypoint.sh && \
    echo 'echo "Waiting for MySQL..."' >> /app/docker-entrypoint.sh && \
    echo 'while ! nc -z mysql 3306; do' >> /app/docker-entrypoint.sh && \
    echo '  sleep 1' >> /app/docker-entrypoint.sh && \
    echo 'done' >> /app/docker-entrypoint.sh && \
    echo 'echo "MySQL started"' >> /app/docker-entrypoint.sh && \
    echo 'LOG_FILE="/app/logs/app.log"' >> /app/docker-entrypoint.sh && \
    echo 'touch $LOG_FILE' >> /app/docker-entrypoint.sh && \
    echo 'echo "Starting application..."' >> /app/docker-entrypoint.sh && \
    echo 'exec /app/apiserver 2>&1 | tee -a $LOG_FILE' >> /app/docker-entrypoint.sh && \
    chmod +x /app/docker-entrypoint.sh

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8091

# Command to run the application
CMD ["/app/docker-entrypoint.sh"]