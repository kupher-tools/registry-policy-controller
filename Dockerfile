# Stage 1:
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN echo "Listing files in /app" && ls -l

COPY cmd internal go.mod ./

RUN echo "Listing files in /app" && ls -l /app

RUN go mod tidy



RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  ./cmd/main.go

# Stage 2:
FROM alpine:3.18

RUN adduser -D webhook

COPY --from=builder /app/main /usr/local/bin/main

USER webhook
WORKDIR /home/webhook

EXPOSE 8443

ENTRYPOINT ["/usr/local/bin/main"]