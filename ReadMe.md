![banner](assets/banner.png)

# ⚽ liga — учёт матчей и рейтингов

Кратко
`liga` — лёгкий и расширяемый REST-сервис для учёта матчей, игроков и таблиц сезона. Поддерживает подсчёт рейтингов (ELO), seed-данные и API для интеграции с фронтендом или админкой.

Getting Started

Prerequisites
- Go >= 1.18
- Git
- (Опционально) Docker, Make

Installation
```bash
git clone https://github.com/your-org/liga.git
cd liga
go mod download
```

Running
- Локально (быстро):
```bash
go run cmd/shumnaya/main.go
```
- Сборка и запуск бинарника:
```bash
go build -o bin/liga ./cmd/shumnaya
./bin/liga
```
- Seed-данные:
```bash
go run cmd/seed/main.go
```

Notes
- Точки входа: `cmd/shumnaya/main.go`, `cmd/seed/main.go`.
- Конфигурация в `internal/config` — проверьте env-переменные перед запуском.

Additional commands
- Тесты:
```bash
go test ./...
```
- Форматирование и линтинг:
```bash
gofmt -w .
golangci-lint run
```
- Docker (если есть Dockerfile):
```bash
docker build -t liga:latest .
docker run -p 8080:8080 liga:latest
```

Contributors
- 👨‍💻 Adam Gowz — основной автор
- 🤝 PRs приветствуются

Ключевая технология
- Язык: Go — надёжный выбор для серверных приложений и микросервисов.
- JWT: `utils/jwt.go` — аутентификация и авторизация.
- Рейтинг: `utils/elo/elo.go` — расчёт ELO для матчей.

Feedback
- Открывайте issues в репозитории.
- Email: developer@your-email.example (замените на реальный).

Спасибо за интерес к проекту!

```mermaid
flowchart LR
  A[API] --> B[Handlers]
  B --> C[Services]
  C --> D[Repositories]
  C --> E[Utils (ELO, JWT)]
  D --> F[(Database)]
```
