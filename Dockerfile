FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./cmd/api

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/api-server .

COPY --from=builder /app/migrations ./migrations

# директорию для логов
RUN mkdir -p /app/logs

RUN chmod +x /app/api-server

EXPOSE 8080 3000 9000

CMD ["./api-server"]