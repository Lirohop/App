FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates bash build-base

WORKDIR /app

# Кэширование зависимостей
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o /app/subscription_service ./cmd/app/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/subscription_service .
COPY config/config.yaml ./config.yaml

EXPOSE 8080
CMD ["./subscription_service"]
