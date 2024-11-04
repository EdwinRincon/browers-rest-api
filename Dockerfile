# Build stage
FROM golang:1.21-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app -ldflags="-w -s" .

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Set environment variables
ENV GIN_MODE=release
ENV TZ=Europe/Madrid

# Expose the application port
EXPOSE 5050

# Run the application
CMD ["./app"]
