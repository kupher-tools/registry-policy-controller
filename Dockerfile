# Start with a minimal Go base image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download Go dependencies
RUN go mod tidy

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  cmd/main.go

# Create small image for deployment
FROM alpine:3.18

# Create a non-root user for security
RUN adduser -D webhook

# Copy TLS certs and binary
COPY --from=builder /app/main /usr/local/bin/main

# Set working directory and permissions
USER webhook
WORKDIR /home/webhook

# Expose the port admission webhook will use
EXPOSE 8443

# Run the webhook server
ENTRYPOINT ["/usr/local/bin/main"]