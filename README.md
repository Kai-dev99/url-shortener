# URL shortener

Сервис создаёт короткие ссылки, перенаправляет по алиасу и поддерживает изменение и удаление ссылок.

## Быстрый запуск

Нужен только Docker с Docker Compose.

git clone <repository-url>
cd url-shortener
cp .env.example .env
docker compose up


# Создать короткую ссылку; alias можно не передавать — он будет сгенерирован сервисом.
curl -u dev:local_dev_only_change_me \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","alias":"example"}' \
  http://localhost:8082/url/

# Перенаправить по короткой ссылке.
curl -i http://localhost:8082/example

# Изменить и удалить ссылку.
curl -X PUT -u dev:local_dev_only_change_me \
  -H "Content-Type: application/json" \
  -d '{"alias":"example","url":"https://example.org"}' \
  http://localhost:8082/url/
curl -X DELETE -u dev:local_dev_only_change_me http://localhost:8082/url/example
```

Для остановки контейнеров используйте `docker compose down`. Добавьте `-v`, только если также нужно удалить данные PostgreSQL.
