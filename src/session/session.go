package session

import (
	"context"
	"sync"
)

type Session struct {
	UserID string
	State  SessionState
	Memory map[string]any
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

func (r *InMemorySessionRepo) GetOrCreate(_ context.Context, id string) *Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s, ok := r.store[id]; ok {
		return s
	}
	s := &Session{UserID: id, State: IdleState{}, Memory: make(map[string]any)}
	r.store[id] = s
	return s
}

func (r *InMemorySessionRepo) Save(_ context.Context, session *Session) error {
	return nil
}
