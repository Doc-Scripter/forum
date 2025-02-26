# Build stage
FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Tidy up Go modules
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o forum .

# Final stage
FROM alpine:3.18

# Install necessary runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    sqlite-libs

# Create a non-root user
RUN adduser -D appuser

# Set the working directory
WORKDIR /app

# Create directories and set permissions
RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

# Copy the built application from the builder stage
COPY --from=builder /app/forum /app/
COPY --from=builder /app/web /app/web

# Set permissions
RUN chown -R appuser:appuser /app && \
    chmod +x /app/forum

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 33333

# Command to run the application
CMD ["/app/forum"]