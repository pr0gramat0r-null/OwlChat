# OwlChat — технічна стратегія (Android / iOS / Windows)

## 0) Аналіз ризиків (перед вибором стеку)

1. **Безпека/E2E**
   - Ризик: помилка в криптографічному протоколі => компрометація приватних чатів.
   - Мітігація: використання battle-tested протоколу **Signal Protocol** (через офіційні/відомі імплементації), мінімум самописної криптографії.

2. **Низький пінг та real-time**
   - Ризик: висока затримка при fan-out повідомлень у великих групах.
   - Мітігація: WebSocket gateway + pub/sub брокер (NATS/Redis Streams/Kafka залежно від масштабу), оптимістичне оновлення UI.

3. **Масштабованість**
   - Ризик: передчасний перехід до мікросервісів у pet-проєкті збільшує складність DevOps.
   - Мітігація: модульний моноліт на старті + чіткі bounded contexts для майбутнього виділення сервісів.

4. **Кросплатформність**
   - Ризик: різна поведінка push, фонових задач, файлового доступу на Android/iOS/Windows.
   - Мітігація: один UI-стек + нативні адаптери для push/secure storage.

5. **Медіафайли**
   - Ризик: висока вартість трафіку/зберігання, повільні завантаження.
   - Мітігація: object storage + CDN + multipart upload + pre-signed URLs.

---

## 1) Вибір стеку (з обґрунтуванням)

### 1.1 Mobile/Desktop

**Рекомендація:** **Flutter (Dart)** для Android/iOS/Windows.

- Плюси:
  - Один кодбейс для 3 платформ.
  - Висока швидкість UI, добра продуктивність.
  - Зріла екосистема для WebSocket, локальної БД, DI, state-management.
- Мінуси:
  - Частину платформених фіч (push, background constraints) доведеться закривати platform channels.

**Альтернатива (якщо команда сильна в TS):** React Native + Electron/Tauri для Windows.
- Мінус: 2 різні runtime-шари для mobile vs desktop (вища підтримка).

### 1.2 Backend

**Рекомендація:** **Go (Golang) + Fiber/Chi + gRPC (internal) + REST/WebSocket (external)**.

- Плюси:
  - Низький latency, ефективна конкурентність.
  - Простий деплой (single binary).
  - Добрий fit для real-time та мережевих навантажень.
- Мінуси:
  - Потреба дисципліни в архітектурі (щоб не виріс «god package»).

### 1.3 Database

**Рекомендація:** **PostgreSQL** (основна), **Redis** (кеш/ephemeral state).

- PostgreSQL:
  - Транзакції, надійність, JSONB для гнучких полів.
  - Повнотекст/індекси для пошуку по повідомленнях.
- Redis:
  - Online-presence, typing, rate-limiting counters, короткоживучі сесії.

### 1.4 Real-time

- Client ↔ Gateway: **WebSocket**.
- Межсервісно (або модулі):
  - MVP: Redis Pub/Sub.
  - Далі при рості: NATS/Kafka (durable streams, replay).

### 1.5 Auth

**Рекомендація:**
- Access token: **JWT (короткий TTL, 10-15 хв)**.
- Refresh token: ротація + зберігання хешу в БД.
- 2FA (TOTP) — опційно на кроці 3.
- Device binding (device_id, fingerprint-lite) для контролю сесій.

---

## 2) Архітектура

### 2.1 Логічна ERD-схема

```text
users
- id (uuid, pk)
- username (unique)
- phone (unique, nullable)
- email (unique, nullable)
- password_hash (nullable, якщо OTP-only)
- public_key_identity (bytea)         -- для E2E
- created_at, updated_at

user_devices
- id (uuid, pk)
- user_id (fk -> users.id)
- device_id (unique per user)
- platform (android|ios|windows)
- push_token
- last_seen_at
- created_at

chats
- id (uuid, pk)
- type (private|group|channel)
- title (nullable)
- created_by (fk -> users.id)
- created_at, updated_at

chat_members
- chat_id (fk -> chats.id)
- user_id (fk -> users.id)
- role (owner|admin|member)
- joined_at
- muted_until (nullable)
- PRIMARY KEY (chat_id, user_id)

messages
- id (uuid, pk)
- chat_id (fk -> chats.id, indexed)
- sender_id (fk -> users.id)
- client_msg_id (string, unique per sender) -- idempotency
- msg_type (text|image|video|file|system)
- body_ciphertext (bytea/text)               -- E2E payload
- body_preview (text, nullable)              -- тільки де дозволено
- reply_to_message_id (nullable fk -> messages.id)
- sent_at (indexed)
- edited_at (nullable)
- deleted_at (nullable)

message_receipts
- message_id (fk -> messages.id)
- user_id (fk -> users.id)
- status (sent|delivered|read)
- ts
- PRIMARY KEY (message_id, user_id)

media_files
- id (uuid, pk)
- uploader_id (fk -> users.id)
- storage_key (unique)        -- шлях в object storage
- mime_type
- size_bytes
- checksum_sha256
- width, height, duration_sec (nullable)
- created_at

message_media
- message_id (fk -> messages.id)
- media_id (fk -> media_files.id)
- ord
- PRIMARY KEY (message_id, media_id)
```

### 2.2 Моноліт vs мікросервіси (для pet-project)

**Рекомендація для старту:** **модульний моноліт**.

- Чому:
  - Менше DevOps-оверхеду (CI/CD, observability, локальний запуск).
  - Швидше ітерації MVP.
  - Легше гарантувати консистентність транзакцій (messages, receipts).
- Як підготуватися до росту:
  - Чітко розділити модулі: `auth`, `chat`, `message`, `media`, `presence`, `admin`.
  - Використовувати доменні інтерфейси + події (outbox pattern).
  - Пізніше винести найгарячіші модулі: `realtime-gateway`, `media-service`.

### 2.3 Безпека

1. **E2E для приватних чатів**
   - Протокол: X3DH + Double Ratchet (Signal-підхід).
   - Сервер зберігає лише ciphertext + метадані доставки.
   - Prekeys на сервері (one-time prekeys), ротація ключів.
   - Для multi-device: окремі сесії ключів на device-level.

2. **Валідація вхідних даних**
   - Схемна валідація DTO (length, enum, format, size limits).
   - Санітизація user-generated content (залежно від формату рендеру).
   - Rate limiting на endpoint/socket-event рівні.
   - Захист від replay: `client_msg_id` + timestamp window.

3. **Інші базові контролі**
   - TLS 1.2+ всюди.
   - Argon2id/bcrypt для password hash.
   - Audit log для admin-дій.

---

## 3) Roadmap реалізації

### Крок 1 — MVP
- Реєстрація/логін (JWT + refresh).
- Створення приватного/групового чату.
- Надсилання/отримання текстових повідомлень через WebSocket.
- Історія повідомлень (pagination cursor-based).
- Мінімальний E2E для private chat (якщо не встигаєте: feature flag і staged rollout).

**Definition of Done:**
- P95 send->receive < 300ms в одній геозоні.
- Немає дублів повідомлень (idempotency).

### Крок 2 — Media + presence
- Upload/download медіа через pre-signed URL.
- Typing status, online/offline presence (Redis TTL).
- Read receipts (delivered/read).

**DoD:**
- Стабільний upload великих файлів (resumable/multipart).
- Presence latency < 3s.

### Крок 3 — Admin + analytics
- Адмін-панель: users/chats moderation, блокування, audit trail.
- Базова аналітика: DAU/MAU, retention D1/D7, message throughput.
- Alerting (error rate, websocket disconnect spikes).

### Крок 4 — Performance
- Кешування hot-даних (chat list, профілі) в Redis.
- CDN для медіа.
- DB-оптимізації: індекси, partitioning по `messages.sent_at` (за потреби).
- Навантажувальні тести (k6/Locust), profiling.

---

## 4) Базова структура проєкту

```text
owlchat/
  apps/
    client_flutter/
      lib/
        core/                 # config, DI, networking, secure storage
        features/
          auth/
          chats/
          messages/
          media/
          settings/
        shared/
      android/
      ios/
      windows/
  services/
    backend/
      cmd/server/
        main.go
      internal/
        app/                  # bootstrap
        modules/
          auth/
          chat/
          message/
          media/
          presence/
          admin/
        platform/
          db/
          cache/
          queue/
          ws/
          storage/
        transport/
          http/
          websocket/
          grpc/
      migrations/
      api/
        openapi.yaml
  deploy/
    docker/
    k8s/
  scripts/
  docs/
    architecture.md
    api.md
```

---

## 5) Приклад ключового модуля (Go, message send)

```go
// internal/modules/message/service.go
package message

import (
    "context"
    "errors"
    "time"
)

type Repository interface {
    ChatMember(ctx context.Context, chatID, userID string) (bool, error)
    ExistsClientMsgID(ctx context.Context, senderID, clientMsgID string) (bool, error)
    InsertMessage(ctx context.Context, m Message) (Message, error)
}

type EventBus interface {
    PublishMessageCreated(ctx context.Context, evt MessageCreatedEvent) error
}

type Service struct {
    repo Repository
    bus  EventBus
}

type SendMessageInput struct {
    ChatID       string
    SenderID     string
    ClientMsgID  string // idempotency key from client
    Ciphertext   []byte // E2E payload
    MessageType  string
    SentAtClient time.Time
}

type Message struct {
    ID          string
    ChatID      string
    SenderID    string
    ClientMsgID string
    Ciphertext  []byte
    MessageType string
    SentAt      time.Time
}

type MessageCreatedEvent struct {
    MessageID string
    ChatID    string
    SenderID  string
    SentAt    time.Time
}

func (s *Service) SendMessage(ctx context.Context, in SendMessageInput) (Message, error) {
    if in.ChatID == "" || in.SenderID == "" || in.ClientMsgID == "" {
        return Message{}, errors.New("invalid input")
    }
    if len(in.Ciphertext) == 0 || len(in.Ciphertext) > 256*1024 {
        return Message{}, errors.New("ciphertext size out of bounds")
    }

    member, err := s.repo.ChatMember(ctx, in.ChatID, in.SenderID)
    if err != nil {
        return Message{}, err
    }
    if !member {
        return Message{}, errors.New("forbidden")
    }

    exists, err := s.repo.ExistsClientMsgID(ctx, in.SenderID, in.ClientMsgID)
    if err != nil {
        return Message{}, err
    }
    if exists {
        return Message{}, errors.New("duplicate client_msg_id")
    }

    msg := Message{
        ChatID:      in.ChatID,
        SenderID:    in.SenderID,
        ClientMsgID: in.ClientMsgID,
        Ciphertext:  in.Ciphertext,
        MessageType: in.MessageType,
        SentAt:      time.Now().UTC(),
    }

    created, err := s.repo.InsertMessage(ctx, msg)
    if err != nil {
        return Message{}, err
    }

    _ = s.bus.PublishMessageCreated(ctx, MessageCreatedEvent{
        MessageID: created.ID,
        ChatID:    created.ChatID,
        SenderID:  created.SenderID,
        SentAt:    created.SentAt,
    })

    return created, nil
}
```

### Потенційні вузькі місця
- **Hot chats**: тисячі повідомлень/сек в один чат → потрібен sharding fan-out та батчинг розсилки.
- **N+1 запити** при завантаженні списку чатів + last message → денормалізація `chat_last_message`.
- **Повторні доставки** при reconnect WebSocket → ack-механізм + idempotent handlers.
- **Media egress cost** → агресивний CDN cache-control + thumbnails/transcoding pipeline.
