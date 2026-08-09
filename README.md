# Project 3: Async — Инфраструктура

## О проекте

В этом проекте вы добавите **асинхронную обработку** через очереди сообщений и объектное хранилище.

**Что вы реализуете:**
- **Worker** — фоновый сервис для обработки изображений
- Публикация и потребление сообщений через RabbitMQ
- Загрузка и хранение файлов в MinIO (S3-совместимое хранилище)
- Асинхронный flow: Upload → Queue → Worker → Storage

**Архитектура:**
```
                    ┌─────────────────┐
Frontend (:5175) ──►│  auth-service   │──► auth-postgres (:5432)
                    │     (:8081)     │
                    └─────────────────┘
                    ┌─────────────────┐
                 ──►│  books-service  │──► books-postgres (:5433)
                    │     (:8082)     │──► MinIO (:9000)
                    └────────┬────────┘
                             │ publish
                             ▼
                    ┌─────────────────┐
                    │    RabbitMQ     │
                    │  (:5672/15672)  │
                    └────────┬────────┘
                             │ consume
                             ▼
                    ┌─────────────────┐
                    │     Worker      │──► MinIO (:9000)
                    └─────────────────┘──► books-postgres (:5433)
```

## Содержимое

```
├── frontend/               # React-приложение (готовое)
├── auth-service/
│   └── migrations/         # Миграции для auth БД
├── books-service/
│   └── migrations/         # Миграции для books БД (включая covers)
├── docker-compose.yml      # PostgreSQL×2 + RabbitMQ + MinIO + Frontend
└── README.md               # Этот файл
```

## Запуск

```bash
docker compose up -d --build
```

После запуска:
- **Frontend**: http://localhost:5175
- **auth-postgres**: localhost:5432
- **books-postgres**: localhost:5433
- **RabbitMQ UI**: http://localhost:15672 (guest/guest)
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)

## Как это работает

Frontend показывает новые возможности по мере реализации:

1. Реализовали Worker → появилась загрузка обложек
2. Загрузили обложку → статус "Processing..."
3. Worker обработал → обложка появляется!

При остановке Worker сообщения сохраняются в очереди RabbitMQ и будут обработаны после перезапуска.

## Подключение к сервисам

**auth-postgres:**
```
postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
```

**books-postgres:**
```
postgres://postgres:postgres@localhost:5433/books?sslmode=disable
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

### Frontend не видит сервисы
1. Убедитесь, что auth-service запущен на порту 8081
2. Убедитесь, что books-service запущен на порту 8082
3. Проверьте CORS middleware в ваших сервисах
4. Откройте DevTools → Network для диагностики

### Ошибка подключения к БД
```bash
docker compose ps                  # Статус контейнеров
docker compose logs auth-postgres  # Логи auth БД
docker compose logs books-postgres # Логи books БД
```

### RabbitMQ не запускается
```bash
docker compose logs rabbitmq  # Логи RabbitMQ
```

### MinIO не запускается
```bash
docker compose logs minio  # Логи MinIO
```

### Пересборка frontend
```bash
docker compose build --no-cache frontend
docker compose up -d frontend
```

## Инструкции по реализации

Подробное описание каждого этапа находится на сайте курса:

**https://praxiscode.io**
