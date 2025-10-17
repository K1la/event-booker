## EventBooker

Легкий сервис для создания мероприятий и управления бронированиями. Проект решает задачу простой записи пользователей на событие с подтверждением оплаты через очередь сообщений. Сервис предоставляет REST API для:
- создания мероприятия,
- бронирования мест на мероприятие,
- подтверждения бронирования (например, после успешной оплаты),
- получения одного мероприятия и списка всех мероприятий.

Бэкенд написан на Go, хранит данные в PostgreSQL и использует RabbitMQ для асинхронной обработки подтверждений. В составе репозитория есть минимальный фронтенд (`web/`) для ручной проверки.

---

## API
Базовый URL: `/api/events`

### POST /api/events
Создать новое мероприятие.
- Тело (JSON):
```json
{
  "title": "Go Meetup",
  "event_at": "2025-10-20T18:00:00Z",
  "total_seats": 100
}
```
- Ответ 200 OK:
```json
{ "result": { /* объект события */ } }
```

### POST /api/events/{id}/book
Забронировать места на мероприятие.
- Тело (JSON):
```json
{
  "telegram_id": 123456789,
  "places_count": 2
}
```
- Ответ 200 OK:
```json
{ "result": { /* объект брони */ } }
```

### POST /api/events/{id}/confirm
Подтвердить бронирование (симулирует успешную оплату). Отправляет сообщение в RabbitMQ и меняет статус брони.
- Тело: пустое
- Ответ 200 OK:
```json
{ "result": { "status": "payment confirmed" } }
```

### GET /api/events/{id}
Получить мероприятие по ID (включая брони, если реализовано на уровне модели/репозитория).
- Ответ 200 OK:
```json
{ "result": { /* объект события */ } }
```

### GET /api/events
Получить список всех мероприятий.
- Ответ 200 OK:
```json
{ "result": [ /* массив событий */ ] }
```

Примечание по ошибкам: при ошибках сервис возвращает
```json
{ "error": "сообщение_об_ошибке" }
```

---

## Запуск через Docker Compose

### 1) Склонировать репозиторий и перейти в каталог проекта
```bash
git clone <repo_url>
cd event-booker
```

### 2) Создать файл .env
Ниже пример значений, подходящих для локального запуска (совместим с `docker-compose.yml`):
```env
# PostgreSQL
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=event-booker

# RabbitMQ
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672

# Goose (миграции)
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=/migrations
```

Важно: файл `env/config.yaml` уже содержит дефолты (`host: db`, `port: 5432` и т.п.), но переменные окружения из `.env` переопределят их при работе контейнеров.

### 3) Собрать контейнеры
```bash
docker compose build
```

### 4) Поднять инфраструктуру
```bash
docker compose up -d
```
Сервисы:
- `db` — PostgreSQL (порт 5432)
- `migrator` — применение миграций через `goose`
- `rabbitmq` — брокер сообщений (порты 5672/AMQP, 15672/менеджмент)
- `backend` — бинарь Go (`./event-booker`) на порту 8080

Проверка логов:
```bash
docker compose logs -f backend
```

### 5) Доступ
- API: `http://localhost:8080/api/events`
- Мини-фронтенд: `http://localhost:8080/` (статические файлы из `web/` раздаются как fallback маршрутом)
- RabbitMQ Management: `http://localhost:15672` (логин/пароль: `guest/guest`)

### 6) Остановка
```bash
docker compose down
```

---

## Дерево проекта
Сокращенная структура каталогов:
```
.
├── cmd
│   └── event-booker
│       └── main.go
├── env
│   └── config.yaml
├── internal
│   ├── api
│   │   ├── handler
│   │   │   ├── get.go
│   │   │   ├── handler.go
│   │   │   ├── interface.go
│   │   │   └── post.go
│   │   ├── response
│   │   │   └── response.go
│   │   ├── router
│   │   │   └── router.go
│   │   └── server
│   │       └── server.go
│   ├── config
│   │   ├── config.go
│   │   └── types.go
│   ├── dto
│   │   └── dto.go
│   ├── model
│   │   └── model.go
│   ├── rabbitmq
│   │   └── rabbitmq.go
│   ├── repository
│   │   ├── create.go
│   │   ├── get.go
│   │   ├── repo.go
│   │   └── update.go
│   ├── sender
│   │   └── sender.go
│   └── service
│       ├── create.go
│       ├── get.go
│       ├── interface.go
│       ├── queue.go
│       ├── service.go
│       └── update.go
├── migrations
│   └── 20251013153247_message_table.sql
├── web
│   ├── index.html
│   ├── admin.html
│   ├── user.html
│   ├── api.js
│   └── styles.css
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Пример запроса и ответа
Пример: создание мероприятия.

HTTP
```http
POST /api/events HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "title": "Go Meetup",
  "event_at": "2025-10-20T18:00:00Z",
  "total_seats": 100
}
```

Успешный ответ
```json
{
  "result": {
    "id": "c0c3f1b5-5e8a-4d1c-9a6f-9c6e3e6d9a11",
    "title": "Go Meetup",
    "event_at": "2025-10-20T18:00:00Z",
    "total_seats": 100,
    "available_seats": 100,
    "bookings": []
  }
}
```

Ошибка (пример)
```json
{ "error": "invalid request body" }
```


