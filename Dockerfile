FROM golang:1.20-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go mod tidy


RUN CGO_ENABLED=1 go build -o forum .

# Use a minimal Alpine image for the final stage
FROM alpine:3.18

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /root/


COPY --from=builder /app/forum .

EXPOSE 33333

CMD ["./forum"]