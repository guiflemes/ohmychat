package ohmychat

import (
	"context"
	"sync"
	"time"
)

const SessionExpiresAt = time.Duration(5) * time.Minute

type Session struct {
	UserID         string
	StateID        StateID
	State          SessionState
	Memory         map[string]any
	LastActivityAt time.Time
}

func (s *Session) IsExpired(timeout time.Duration) bool {
	return time.Since(s.LastActivityAt) > timeout
}

type InMemorySessionRepo struct {
	mu    sync.Mutex
	store map[string]*Session
}

func NewInMemorySessionRepo() *InMemorySessionRepo {
	return &InMemorySessionRepo{
		store: make(map[string]*Session),
	}
}

func (r *InMemorySessionRepo) GetOrCreate(_ context.Context, id string) (*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s, ok := r.store[id]; ok {
		return s, nil
	}
	s := &Session{UserID: id, State: IdleState{}, Memory: make(map[string]any), LastActivityAt: time.Now()}
	r.store[id] = s
	return s, nil
}

func (r *InMemorySessionRepo) Save(_ context.Context, session *Session) error {
	return nil
}
