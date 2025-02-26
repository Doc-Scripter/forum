# Build stage
FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory
WORKDIR /app


# Copy the rest of the source code
COPY . .

# Tidy up Go modules
RUN go mod tidy


# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o forum .

# Final stage
FROM alpine:3.18

# Install necessary runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    sqlite-libs


# Set the working directory
WORKDIR /app

# Create directories and set permissions
RUN mkdir -p /app/data

# Copy the built application from the builder stage
COPY --from=builder /app/forum /app/

COPY --from=builder /app/web /app/web


# Expose the application port
EXPOSE 33333

# Command to run the application
CMD ["/app/forum"]