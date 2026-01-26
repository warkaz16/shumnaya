![banner](assets/banner.png)

# liga — учёт матчей и рейтингов

## Описание
`liga` — лёгкий и расширяемый REST-сервис для учёта матчей, игроков и таблиц сезона. Поддерживает подсчёт рейтингов (ELO), seed-данные и API для интеграции с фронтендом или админкой.

## Cтек
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)]

## Prerequisites
- Go >= 1.18
- Git
- (Опционально) Docker, Make

## Архитектура
```mermaid
flowchart LR
  A[API] --> B[Handlers]
  B --> C[Services]
  C --> D[Repositories]
  C --> E[Utils — ELO & JWT]
  D --> F[(Database)]
```

## Примеры запросов
```http
GET /matches?from=2025-01-01&to=2025-12-31
GET /matches?season_id=1&from=2024-12-25&to=2025-01-25
GET /matches?player_id=5&from=2024-12-20&to=2024-12-25
```

## Установка
```bash
git clone https://github.com/your-org/liga.git
cd liga
go mod download
```

## Запуск локально
- Быстро (локально):
```bash
go run cmd/liga/main.go
```
- Сборка и запуск бинарника:
```bash
go build -o bin/liga ./cmd/liga
./bin/liga
```
- Запустить seed-данные:
```bash
go run cmd/seed/main.go
```

## Примечания
- Точки входа: [cmd/liga/main.go](cmd/liga/main.go), [cmd/seed/main.go](cmd/seed/main.go).
- Конфигурация: смотрите [internal/config/database.go](internal/config/database.go).

## Дополнительные команды
- Тесты:
```bash
go test ./...
```
- Форматирование и линт:
```bash
gofmt -w .
golangci-lint run
```
- Docker (если есть Dockerfile):
```bash
docker build -t liga:latest .
docker run -p 8080:8080 liga:latest
```

## Участники
- 👨‍💻 https://github.com/warkaz16 — основной автор
- 🤝 PRs и issues приветствуются

## Ключевая технология
- JWT: реализовано в [utils/jwt.go](utils/jwt.go).
- ELO: расчёт в [utils/elo/elo.go](utils/elo/elo.go).

## Обратная связь
- Откройте issue в репозитории или пишите на warkaz16@mail.ru