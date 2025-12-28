# REST API task manager written mainly on Go

## Запуск локально

### 1) Клонируйте репозиторий
``
git clone https://github.com/maket12/todo-api
``

``
cd todo-api
``

### 2) Запуск приложения
``
go run cmd/todo/main.go
``

Сервер будет доступен по адресу: http://localhost:8080.

---

## Запуск в Docker
### 1) Сборка образа
``
docker build -t todo-api .
``

### 2) Запуск контейнера
``
docker run -p 8080:8080 todo-api
``
