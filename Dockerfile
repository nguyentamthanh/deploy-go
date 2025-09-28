# Build stage
FROM golang:1.21-alpine3.18 AS builder

WORKDIR /app

# Install git and ca-certificates for go mod download
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a -installsuffix cgo \
  -o main .

# Final stage
FROM alpine:3.20

# Install ca-certificates for HTTPS requests and dumb-init
RUN apk --no-cache add ca-certificates tzdata dumb-init

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
  adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser:appgroup /app/main && \
  chmod 755 /app/main

# Switch to non-root user
USER appuser

# Set secure environment variables
ENV GIN_MODE=release

# Expose port
EXPOSE 8080

# Run the application with dumb-init for proper signal handling
ENTRYPOINT ["dumb-init", "--"]
CMD ["./main"]
