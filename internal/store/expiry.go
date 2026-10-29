package store

import "time"

func (s *Store) SetExpiry(key string , t time.Duration) bool {
	s.mu.Lock()

	defer s.mu.Unlock()

	e, exists := s.data[key]
	if !exists || e.expired(time.Now()) {
		return false
	}

	if t <= 0 {
		delete(s.data, key)
		return true
	}

	e.expiresAt = time.Now().Add(t)

	// entry is stored by value, so write the updated copy back
	s.data[key] = e
	return true
}
