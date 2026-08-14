# Project 4: API Gateway — Инфраструктура

## О проекте

В этом проекте вы создадите **API Gateway** — единую точку входа для всех клиентов.

**Что вы реализуете:**
- **api-gateway** — reverse proxy с маршрутизацией
- Rate limiting на базе Redis
- Кэширование запросов
- CORS, логирование, трейсинг
- Dashboard с агрегированными данными

**Архитектура:**
```
                         ┌───────────────────────────────────────┐
                         │           API Gateway (:8000)          │
                         │  ┌─────────┐ ┌─────────┐ ┌─────────┐  │
Frontend (:5176) ───────►│  │  CORS   │→│  Rate   │→│  Cache  │  │
                         │  │         │ │ Limiter │ │         │  │
                         │  └─────────┘ └─────────┘ └─────────┘  │
                         │                   │                    │
                         └───────────────────┼────────────────────┘
                                             │
                    ┌────────────────────────┼────────────────────────┐
                    │                        │                        │
                    ▼                        ▼                        ▼
           ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
           │  auth-service   │    │  books-service  │    │     Worker      │
           │     (:8081)     │    │     (:8082)     │    │                 │
           └────────┬────────┘    └────────┬────────┘    └────────┬────────┘
                    │                      │                      │
                    ▼                      ▼                      ▼
           auth-postgres(:5432)   books-postgres(:5433)   RabbitMQ / MinIO
```

## Содержимое

```
├── frontend/               # React-приложение (готовое)
├── auth-service/
│   └── migrations/         # Миграции для auth БД
├── books-service/
│   └── migrations/         # Миграции для books БД
├── docker-compose.yml      # PostgreSQL×2 + Redis + RabbitMQ + MinIO + Frontend
└── README.md               # Этот файл
```

## Запуск

```bash
docker compose up -d --build
```

После запуска:
- **Frontend**: http://localhost:5176
- **auth-postgres**: localhost:5432
- **books-postgres**: localhost:5433
- **Redis**: localhost:6379
- **RabbitMQ UI**: http://localhost:15672 (guest/guest)
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)

## Как это работает

Frontend настроен на работу через ваш gateway (порт 8000). По мере реализации функций gateway — они становятся доступны:

1. Реализовали proxy → запросы проходят через gateway
2. Добавили CORS → frontend перестаёт получать ошибки
3. Добавили rate limiting → защита от перегрузки
4. Добавили кэш → ускорение повторных запросов

## Подключение к сервисам

**auth-postgres:**
```
postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
```

**books-postgres:**
```
postgres://postgres:postgres@localhost:5433/books?sslmode=disable
```

**Redis:**
```
redis://localhost:6379
```

**RabbitMQ:**
```
amqp://guest:guest@localhost:5672/
```

**MinIO:**
```
Endpoint: localhost:9000
Access Key: minioadmin
Secret Key: minioadmin
```

## Тестовые пользователи

| Email | Password |
|-------|----------|
| admin@bookshelf.dev | password123 |
| john@example.com | password123 |
| maria@example.com | password123 |

## Команды

```bash
docker compose up -d --build  # Запустить
docker compose down           # Остановить
docker compose logs -f        # Логи
docker compose down -v        # Удалить всё (включая данные)
```

## Устранение неполадок

### Docker не запускается
```bash
docker ps  # Проверьте, что Docker daemon работает
```

### Frontend не видит gateway
1. Убедитесь, что api-gateway запущен на порту 8000
2. Проверьте CORS middleware в gateway
3. Откройте DevTools → Network для диагностики

### Ошибка 502 Bad Gateway
Backend-сервисы не запущены или недоступны:
1. Убедитесь, что auth-service запущен на порту 8081
2. Убедитесь, что books-service запущен на порту 8082
3. Проверьте логи сервисов

### Ошибка подключения к БД
```bash
docker compose ps                  # Статус контейнеров
docker compose logs auth-postgres  # Логи auth БД
docker compose logs books-postgres # Логи books БД
```

### Redis не запускается
```bash
docker compose logs redis  # Логи Redis
```

### Пересборка frontend
```bash
docker compose build --no-cache frontend
docker compose up -d frontend
```

## Инструкции по реализации

Подробное описание каждого этапа находится на сайте курса:

**https://praxiscode.io**
