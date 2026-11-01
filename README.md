# QR Code Generator

Сервис создаёт QR-коды для HTTP- и HTTPS-ссылок. Готовое PNG-изображение можно сразу забрать из ответа API или позже получить по UUID. Для хранения используется Redis с ограниченным сроком жизни записей.

## Запуск

Нужны Docker Engine и Docker Compose v2. При первом запуске подготовь локальные файлы конфигурации:

```bash
cp backend/.env.example backend/.env
cp backend/config/config.example.yaml backend/config/config.yaml
```

Для запуска через Compose укажи `redis:6379` в поле `database.host` файла `backend/config/config.yaml`: именно так приложение обращается к Redis внутри Docker-сети.

Затем из корня проекта запусти:

```bash
docker compose up --build -d
```

Интерфейс будет доступен на http://localhost:3000, API — на http://localhost:8080.

Полезные команды:

```bash
docker compose ps
docker compose logs -f app
docker compose logs -f redis
docker compose down
```

`docker compose down` оставляет данные Redis в Docker volume. Чтобы удалить и их, выполни `docker compose down -v` — это действие необратимо.

## API

Полная спецификация OpenAPI 3.0 находится в [`backend/scheme/openapi.yaml`](backend/scheme/openapi.yaml). Её можно открыть в любом инструменте, который поддерживает OpenAPI.

### Создать QR-код

`POST /gen` принимает JSON с абсолютным URL. Поддерживаются только ссылки с протоколом `http` или `https`.

```bash
curl -X POST http://localhost:8080/gen \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

При успехе сервер отвечает `201 Created`:

```json
{
  "id": "8f4c7f61-8c43-4efb-83fd-63b2ae712934",
  "image": "<PNG в Base64>"
}
```

Поле `image` содержит Base64 самого PNG, без префикса `data:image/png;base64,`. Значение `id` понадобится, чтобы получить файл позже.

Если JSON некорректен или ссылка не прошла проверку, API возвращает `400 Bad Request` с телом вида `{"error":"..."}`. Ошибка генерации или сохранения возвращает `500 Internal Server Error` в том же формате.

### Получить изображение

`GET /code/{id}` возвращает PNG с типом содержимого `image/png`:

```bash
curl http://localhost:8080/code/8f4c7f61-8c43-4efb-83fd-63b2ae712934 \
  --output qr.png
```

Если такого UUID нет в Redis или срок хранения записи закончился, сервер отвечает `404 Not Found` и JSON с полем `error`.

### Проверить доступность API

`GET /health` возвращает `200 OK` и `{"status":"ok"}`.

## Конфигурация и разработка

Пример настроек находится в `backend/config/config.example.yaml`, переменных окружения — в `backend/.env.example`. Локальные файлы `backend/config/config.yaml` и `backend/.env` не хранятся в Git. Срок хранения задаётся параметром `database.ttl` в часах; в примере конфигурации он равен 24 часам.

Для запуска Go-проверок из корня проекта:

```bash
make test
make build
```

В Docker Compose Redis доступен приложению по адресу `redis:6379`, а с хоста наружу опубликованы только порты приложения и интерфейса.
