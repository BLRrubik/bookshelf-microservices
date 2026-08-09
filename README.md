# Project 2: Микросервисы — Инфраструктура

## О проекте

В этом проекте вы **декомпозируете монолит** на два независимых микросервиса.

**Что вы реализуете:**
- **auth-service** — авторизация, JWT, профили пользователей
- **books-service** — книги, рецензии, межсервисное взаимодействие
- HTTP-клиент для связи между сервисами
- Graceful shutdown для корректного завершения

**Архитектура:**
```
                    ┌─────────────────┐
Frontend (:5174) ──►│  auth-service   │──► auth-postgres (:5432)
                    │     (:8081)     │
                    └─────────────────┘
                    ┌─────────────────┐
                 ──►│  books-service  │──► books-postgres (:5433)
                    │     (:8082)     │
                    └─────────────────┘
```

## Содержимое архива

В этом архиве только Docker-конфигурация для запуска всей системы:

```
├── docker-compose.yml      # БД (auth-postgres, books-postgres) + auth-service + books-service + frontend
├── auth-service/
│   └── Dockerfile          # Сборка auth-service (Go)
├── books-service/
│   └── Dockerfile          # Сборка books-service (Go)
└── README.md
```

Директории **frontend/** и **migrations/** (внутри auth-service и books-service), а также ваш Go-код сервисов должны уже быть в проекте — из шага с инфраструктурой и миграциями. Положите Dockerfile в соответствующие директории и замените корневой `docker-compose.yml` на файл из архива.

## Запуск

```bash
docker compose up -d --build
```

После запуска:
- **Frontend**: http://localhost:5174
- **auth-service**: http://localhost:8081
- **books-service**: http://localhost:8082
- **auth-postgres**: localhost:5432 (БД `auth`)
- **books-postgres**: localhost:5433 (БД `books`)

## Как это работает

Frontend показывает **статус каждого микросервиса** в футере:

| Индикатор | Значение |
|-----------|----------|
| **Auth ●** (зелёный) | auth-service работает |
| **Books ●** (зелёный) | books-service работает |
| **✕** (красный) | Сервис недоступен |

По мере реализации сервисов — индикаторы будут становиться зелёными, а функции — доступными.

## Подключение к базам данных

**auth-postgres:**
```
postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
```

**books-postgres:**
```
postgres://postgres:postgres@localhost:5433/books?sslmode=disable
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
docker compose ps                # Статус контейнеров
docker compose logs auth-postgres   # Логи auth БД
docker compose logs books-postgres  # Логи books БД
```

### Пересборка frontend
```bash
docker compose build --no-cache frontend
docker compose up -d frontend
```

### Пересборка Go-сервисов (после изменения кода)
```bash
docker compose build --no-cache auth-service books-service
docker compose up -d auth-service books-service
```

## Инструкции по реализации

Подробное описание каждого этапа находится на сайте курса:

**https://praxiscode.io**

## Асинхронная обработка

### Новые компоненты

- **RabbitMQ** — очередь сообщений для асинхронных задач
- **MinIO** — S3-совместимое хранилище для обложек книг
- **worker-service** — фоновый обработчик задач из очереди

### Сценарий обработки обложки

1. Пользователь загружает изображение → books-service сохраняет оригинал в MinIO
2. books-service публикует сообщение в RabbitMQ
3. worker-service забирает сообщение, обрабатывает изображение (resize)
4. worker-service загружает результат в MinIO и обновляет статус в БД

### Почему асинхронно?

Обработка изображений — тяжёлая операция (1-5 секунд). Если делать синхронно,
пользователь будет ждать. Асинхронный подход: пользователь сразу получает ответ
"обложка обрабатывается", а результат появляется в фоне.