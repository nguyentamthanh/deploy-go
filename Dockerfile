# Build stage
FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o main .

# Final stage
FROM alpine:3.20

# Install ca-certificates for HTTPS requests and update packages
RUN apk --no-cache add ca-certificates tzdata && \
  apk --no-cache upgrade && \
  rm -rf /var/cache/apk/* && \
  apk --no-cache add --no-scripts --no-deps dumb-init

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
