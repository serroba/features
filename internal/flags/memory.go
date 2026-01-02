package flags

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	flags map[string]*Flag
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		flags: make(map[string]*Flag),
	}
}

func (r *MemoryRepository) Get(_ context.Context, key string) (*Flag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, ok := r.flags[key]
	if !ok {
		return nil, ErrFlagNotFound
	}

	return flag, nil
}

func (r *MemoryRepository) Create(_ context.Context, flag *Flag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flags[flag.Key]; exists {
		return ErrFlagExists
	}

	r.flags[flag.Key] = flag

	return nil
}
