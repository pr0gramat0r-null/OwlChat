package chat

import (
	"errors"
	"sync"
)

type Chat struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Members []string `json:"members"`
}

type MemoryStore struct {
	mu    sync.RWMutex
	chats map[string]Chat
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{chats: map[string]Chat{}} }

func (s *MemoryStore) Create(c Chat) (Chat, error) {
	if c.ID == "" || len(c.Members) == 0 {
		return Chat{}, errors.New("invalid chat")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chats[c.ID] = c
	return c, nil
}

func (s *MemoryStore) Get(chatID string) (Chat, bool) {
	s.mu.RLock(); defer s.mu.RUnlock()
	c, ok := s.chats[chatID]
	return c, ok
}
