# QR Code Generator

Небольшой сервис для генерации QR-кодов из HTTP/HTTPS URL.

## Состав

- Go + Gin — HTTP API.
- Redis — хранение PNG по UUID с TTL.
- Nginx — статический frontend и reverse proxy к API.
- Docker Compose — запуск всего стенда.

## Запуск

Требуется Docker Engine и Docker Compose v2.

Из корня проекта:

```bash
docker compose up --build -d
```

Открыть интерфейс: http://localhost:3000

API доступно напрямую: http://localhost:8080

Проверка состояния:

```bash
docker compose ps
docker compose logs -f app
docker compose logs -f redis
```

Остановка:

```bash
docker compose down
```

Удаление Redis-данных (необратимо):

```bash
docker compose down -v
```

## API

### Создать QR-код

```http
POST /gen
Content-Type: application/json

{"url":"https://example.com"}
```

Ответ `201 Created`:

```json
{
  "id": "<uuid>",
  "image": "<base64 PNG>"
}
```

### Получить QR-код

```http
GET /code/{id}
```

Успешный ответ — бинарное изображение с `Content-Type: image/png`.

### Health check

```http
GET /health
```

## Проверка через curl

```bash
curl -X POST http://localhost:8080/gen \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

Скопируй `id` из JSON-ответа:

```bash
curl http://localhost:8080/code/<UUID> --output qr.png
```

## Конфигурация

Docker-конфигурация находится в `backend/config/config.yaml`.
Локальные переменные окружения — в `backend/.env`. Не добавляй реальные секреты в Git.

Redis доступен приложению по адресу `redis:6379` внутри Compose-сети. Порт Redis наружу не публикуется.

## Примечания

- PNG сохраняется в Redis как бинарное значение, без Base64.
- Base64 используется только в JSON-ответе `POST /gen`.
- Срок хранения QR-кода задаётся параметром `database.ttl` в часах.
- Redis AOF включён, данные хранятся в Docker volume `redis_data`.
