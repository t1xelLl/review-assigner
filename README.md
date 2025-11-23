# PR Reviewer Assignment Service

Микросервис для автоматического назначения ревьюверов на Pull Request'ы и управления командами.

## Технологии

- Go 1.24.3
- PostgreSQL
- Gin Web Framework
- Docker & Docker Compose
- SQLx
- Migrate

###Установка и запуск

1. **Клонируйте репозиторий**:
```bash
git clone https://github.com/t1xelLl/review-assigner.git
cd review-assigner

```
2. **Установите необходимые пакеты Go**
```bash
go mod tidy
```
3. **Запуск приложения**
```bash
docker compose up -d 
```
