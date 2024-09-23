package inmemory

import (
	"context"
	"sync"

	"github.com/tmvrus/key-value-storage/internal/domain"
)

type engine struct {
	lock sync.RWMutex
	data map[string]string
}

func New() *engine {
	return &engine{data: make(map[string]string)}
}

func (e *engine) Set(_ context.Context, key, value string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.data[key] = value
	return nil
}

func (e *engine) Get(_ context.Context, key string) (string, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	v, ok := e.data[key]
	if !ok {
		return "", domain.ErrNotFound
	}

	return v, nil
}

func (e *engine) Delete(_ context.Context, key string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	_, ok := e.data[key]
	if !ok {
		return domain.ErrNotFound
	}

	delete(e.data, key)
	return nil
}
