package telegram

import "sync"

var MemoryFlowRepo = NewFlowRepo()

type flowRepo struct {
	flows map[string]*Flow
	lock  sync.Mutex
}

func NewFlowRepo() *flowRepo {
	return &flowRepo{
		flows: make(map[string]*Flow),
	}
}

func (r *flowRepo) Save(f *Flow) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.flows[f.Key] = f
	return nil
}

func (r *flowRepo) Get(key string) *Flow {
	r.lock.Lock()
	defer r.lock.Unlock()
	f, ok := r.flows[key]
	if !ok {
		return nil
	}

	return f
}
