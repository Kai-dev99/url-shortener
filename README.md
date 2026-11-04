URL Shortener

Сервис для создания коротких ссылок, перенаправления пользователей по алиасу, а также изменения и удаления существующих ссылок.

Быстрый запуск

Для запуска проекта необходимы Docker и Docker Compose.

# 1. Клонирование репозитория

git clone <repository-url>
cd url-shortener

# 2. Настройка окружения

Создайте файл .env на основе примера:

cp .env.example .env

# 3. Запуск приложения

docker compose up

После запуска сервис будет доступен по адресу:

http://localhost:8082
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
