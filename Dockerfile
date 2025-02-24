# Build stage
  FROM golang:1.20-alpine AS builder

  # Install build dependencies
  RUN apk add --no-cache gcc musl-dev

  # Set the working directory
  WORKDIR /app


  # Copy the source code
  COPY . .

  # Tidy up Go modules
  RUN go mod tidy

  # Build the Go application
  RUN CGO_ENABLED=1 go build -o forum .

  # Final stage
  FROM alpine:3.18

  # Install necessary runtime dependencies
  RUN apk add --no-cache ca-certificates

  # Set the working directory
  WORKDIR /root/

  # Create a directory for the database file
RUN mkdir -p /root/data

  # Copy the built application from the builder stage
  COPY --from=builder /app/forum .

  # Copy the built application template from the builder stage
  COPY --from=builder /app/web /root/web

  # Ensure the forum executable is in the correct location
  RUN chmod +x /root/forum

  # Expose the application port
  EXPOSE 33333

  # Command to run the application
  CMD ["/root/forum"]