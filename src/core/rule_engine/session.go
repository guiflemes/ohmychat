package rule_engine

import (
	"context"
	"sync"
)

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
	s := &Session{UserID: id, State: IdleState{}, Memory: map[string]string{}}
	r.store[id] = s
	return s
}
