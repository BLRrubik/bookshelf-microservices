# Этап 1: сборка
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Git нужен для загрузки некоторых Go-модулей
RUN apk add --no-cache git

# Сначала go.mod — для кэширования слоёв
COPY go.mod go.sum ./
RUN go mod download

# Исходный код
COPY . .

# Статический бинарник для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker-service ./cmd/worker

# Этап 2: рантайм
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /worker-service .

EXPOSE 8083

CMD ["./worker-service"]
