# Multi-stage build for minimal image size
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -o s3-mirror \
  ./cmd/s3-mirror

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/s3-mirror .

# Create a non-root user
RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser && \
  chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["/app/s3-mirror"]
CMD ["-config", "/config/config.yaml"]
