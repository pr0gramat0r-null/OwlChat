# 　OwlChat
　


╔═════════════════════════════════════════════════╗                            
║　　　︾▽︾⩘　　　Messenger pet-project — ***Owl Chat***　　　︾▽︾⩘　　　║  
╚══╔═══════════════════════════════════════════╗══╝  
　　║　⪧ 　　　Cross–platform App for *Android, iOS, Windows*　　　⪦ 　║              
　　╚═══════════════════════════════════════════╝  
  　
## Поточний стан

У репозиторії стартово реалізовано:
- Базовий backend:
  - health-check,
  - dev login,
  - create chats,
  - send/read text message.
- Локальна інфраструктура для dev: PostgreSQL + Redis (`docker-compose`).

> Важливо: backend використовує "in-memory stor", щоб швидко стартувати MVP-розробку.

---

## Структура проєкту

```text
.
├── README.md
├── TECH_STRATEGY.md
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
- Docker + Docker Compose

### 2) Mobile/Desktop
- Flutter SDK
- Android Studio (android SDK)
- Visual Studio (desktop development tools)

---

## Встановлення і запуск

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

### Dev login
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

## Подальші плани

- [ ] Замінити in-memory storage на PostgreSQL + migrations.
- [ ] Додати Redis для presence/typing/rate limiting.
- [ ] Додати WebSocket gateway для реального real-time fan-out.
- [ ] Реалізувати auth middleware (Bearer JWT) та refresh token ротацію.
- [ ] Створити Flutter-клієнт (auth/chats/messages) з clean architecture.
- [ ] Додати e2e encryption для private chats (Signal-підхід) поетапно.

---

