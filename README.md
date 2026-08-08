# Bookshelf — Microservices

## Почему микросервисы?

Монолит стал узким местом

## Архитектура

- **auth-service** (порт 8081) — регистрация, авторизация, управление пользователями
Файлы auth-service определены: 
- domain/user.go
- handler/auth_handler.go
- service/user_service.go
- repository/user_repository.go

- **books-service** (порт 8082) — каталог книг, рецензии 
Файлы books-service определены: 
- domain/book.go
- domain/review.go
- handler/book_handler.go
- handler/review_handler.go
- service/book_service.go
- service/review_service.go
- repository/book_repository.go
- repository/review_repository.go

Каждый сервис имеет свою базу данных (Database per Service).

## Компоненты системы

| Компонент      | Назначение                    |
|----------------|-------------------------------|
| auth-service   | Аутентификация и пользователи |
| books-service  | Книги и рецензии              |
| auth-postgres  | БД для auth-service           |
| books-postgres | БД для books-service          |
| frontend       | React-приложение              |

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

## Содержимое

```
├── frontend/               # React-приложение (готовое)
├── auth-service/
│   └── migrations/         # Миграции для auth БД
├── books-service/
│   └── migrations/         # Миграции для books БД
├── docker-compose.yml      # auth-postgres + books-postgres + Frontend
└── README.md               # Этот файл
```

## Запуск

```bash
docker compose up -d --build
```

После запуска:
- **Frontend**: http://localhost:5174
- **auth-postgres**: localhost:5432
- **books-postgres**: localhost:5433

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

## Инструкции по реализации

Подробное описание каждого этапа находится на сайте курса:

**https://praxiscode.io**
