# OwlChat

Кросплатформний pet-project чат-додаток (Android, iOS, Windows) з фокусом на безпеку, масштабованість і low-latency real-time.

## Поточний стан

У репозиторії стартово реалізовано:
- Архітектурна стратегія: `TECH_STRATEGY_UA.md`.
- Базовий backend (Go) з API для:
  - health-check,
  - dev login (JWT),
  - створення чатів,
  - відправки/читання текстових повідомлень.
- Локальна інфраструктура для dev: PostgreSQL + Redis (`docker-compose`).

> Важливо: поточна реалізація backend використовує in-memory store (без персистентності), щоб швидко стартувати MVP-розробку.

---

## Структура проєкту

```text
.
├── README.md
├── TECH_STRATEGY_UA.md
├── apps
│   └── client_flutter
├── deploy
│   └── docker
│       └── docker-compose.yml
└── services
    └── backend
        ├── .env.example
        ├── Makefile
        ├── go.mod
        ├── cmd/server/main.go
        ├── internal
        │   ├── app
        │   ├── config
        │   ├── modules
        │   │   ├── auth
        │   │   ├── chat
        │   │   └── message
        │   └── platform
        │       └── httpx
        └── migrations
```

---

## Вимоги (Dependencies)

### 1) Backend
- Go `1.22+`
- Docker + Docker Compose (для PostgreSQL/Redis)

### 2) Mobile/Desktop (наступний етап)
- Flutter stable SDK
- Android Studio / Xcode / Visual Studio (Windows desktop workload)

---

## Встановлення і запуск

## 0. Клонування
```bash
git clone <your-repo-url>
cd OwlChat
```

## 1. Підняти локальну інфраструктуру
```bash
docker compose -f deploy/docker/docker-compose.yml up -d
```

Перевірка контейнерів:
```bash
docker compose -f deploy/docker/docker-compose.yml ps
```

## 2. Налаштувати env backend
```bash
cp services/backend/.env.example services/backend/.env
```

За потреби змініть `JWT_SECRET`.

## 3. Встановити Go-залежності
```bash
cd services/backend
go mod tidy
```

## 4. Запустити backend
```bash
make run
```

Сервер стартує на `http://localhost:8080`.

---

## API (MVP scaffold)

### Health
```http
GET /healthz
```

### Dev login (отримати access token)
```http
POST /api/v1/auth/dev-login
Content-Type: application/json

{
  "user_id": "u1"
}
```

### Створити чат
```http
POST /api/v1/chats/
Content-Type: application/json

{
  "id": "chat-1",
  "title": "General",
  "members": ["u1", "u2"]
}
```

### Надіслати повідомлення
```http
POST /api/v1/messages/
Content-Type: application/json

{
  "id": "m1",
  "chat_id": "chat-1",
  "sender_id": "u1",
  "client_msg_id": "u1-0001",
  "body": "Hello"
}
```

### Отримати список повідомлень чату
```http
GET /api/v1/messages/chat-1
```

---

## Команди розробника

З директорії `services/backend`:

```bash
make fmt    # gofmt
make test   # go test ./...
make build  # binary in ./bin
make run    # run server
```

---

## Що робити далі (next milestones)

1. Замінити in-memory storage на PostgreSQL + migrations.
2. Додати Redis для presence/typing/rate limiting.
3. Додати WebSocket gateway для реального real-time fan-out.
4. Реалізувати auth middleware (Bearer JWT) та refresh token ротацію.
5. Створити Flutter-клієнт (auth/chats/messages) з clean architecture.
6. Додати e2e encryption для private chats (Signal-підхід) поетапно.

---

## Нотатки по безпеці

- Поточний `dev-login` endpoint призначений лише для локальної розробки.
- Для production потрібні:
  - повноцінна реєстрація/автентифікація,
  - secure secret management,
  - TLS termination,
  - rate limiting і audit logging.
