# Этап 1: сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Git нужен для загрузки некоторых Go-модулей
RUN apk add --no-cache git

# Сначала go.mod — для кэширования слоёв (go.sum не обязателен)
COPY go.mod ./
RUN go mod download

# Исходный код
COPY . .

# Статический бинарник для Linux (точка входа — cmd/server, как в этапах проекта)
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/server

# Этап 2: рантайм
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /auth-service .

EXPOSE 8081

CMD ["./auth-service"]
