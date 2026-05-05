package message

import "sync"

type Message struct {
	ID string `json:"id"`
	ChatID string `json:"chat_id"`
	SenderID string `json:"sender_id"`
	ClientMsgID string `json:"client_msg_id"`
	Body string `json:"body"`
}

type MemoryStore struct {
	mu sync.RWMutex
	byChat map[string][]Message
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{byChat: map[string][]Message{}} }

func (s *MemoryStore) Add(m Message) Message {
	s.mu.Lock(); defer s.mu.Unlock()
	s.byChat[m.ChatID] = append(s.byChat[m.ChatID], m)
	return m
}

func (s *MemoryStore) List(chatID string) []Message {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]Message, len(s.byChat[chatID]))
	copy(out, s.byChat[chatID])
	return out
}
