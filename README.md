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