# URL Shortener

Сервис сокращения ссылок на Go + Postgres.

## Быстрый старт

```bash
git clone <repo-url> url-shortener
cd url-shortener
cp .env.example .env
docker compose up
```

После старта сервис доступен на [http://localhost:8082](http://localhost:8082).

## Порты и креды

| Что | Значение |
| --- | --- |
| HTTP | `localhost:8082` |
| Basic Auth | `myuser` / `mypass` |
| Postgres | `localhost:5432` |
| Postgres user / password / db | `postgres` / `mypass` / `urlshortener` |
| `DATABASE_URL` | `postgres://postgres:mypass@localhost:5432/urlshortener` |

Эти значения совпадают в `.env.example`, `docker-compose.yml`, `config/local.yaml` и `config/docker.yaml`.

## Примеры запросов

Создать короткую ссылку:

```bash
curl -u myuser:mypass -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"https://example.com\",\"alias\":\"demo\"}"
```

Перейти по короткой ссылке:

```bash
curl -i http://localhost:8082/demo
```

## Локальный запуск без Docker (опционально)

Нужны Go и уже поднятый Postgres с теми же кредами, что выше.

```bash
cp .env.example .env
go run ./cmd/url-shortener
```
