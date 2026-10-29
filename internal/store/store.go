package store

import (
	"sync"
	"time"
)

type Store struct {
	mu sync.Mutex
	data map[string]entry
}

type entry struct {
	value string
	expiresAt time.Time
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]entry),
	}
}

func (e *entry) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && now.After(e.expiresAt)
}

func (s *Store) Set(key string, value string, t time.Duration) {
	s.mu.Lock()

	defer s.mu.Unlock()
	
	e := entry{value: value}

	if t > 0 {
		e.expiresAt = time.Now().Add(t)
    }

	s.data[key] = e
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists {
		return "", false
	}

	// lazy expiry: drop the key on read once its deadline has passed
	if e.expired(time.Now()) {
		delete(s.data, key)
		return "", false
	}

	return e.value, true
}

func (s *Store) Delete(key string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    delete(s.data, key)
}