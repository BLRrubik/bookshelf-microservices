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

# Статический бинарник для Linux (точка входа — cmd/server, как в этапах проекта)
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/server

# Этап 2: рантайм
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /api-gateway .

EXPOSE 8000

CMD ["./api-gateway"]
