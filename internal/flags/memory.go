package flags

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	flags map[FlagKey]Flag
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		flags: make(map[FlagKey]Flag),
	}
}

func (r *MemoryRepository) Get(_ context.Context, key FlagKey) (Flag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flag, ok := r.flags[key]
	if !ok {
		return Flag{}, ErrFlagNotFound
	}

	return flag, nil
}

func (r *MemoryRepository) Create(_ context.Context, flag Flag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flags[flag.Key]; exists {
		return ErrFlagExists
	}

	r.flags[flag.Key] = flag

	return nil
}
