package inmemory

import (
	"context"

	"github.com/tmvrus/key-value-storage/internal/domain"
)

const initialSize = 100

type engine struct {
	data map[string]string
}

func New() *engine {
	return &engine{data: make(map[string]string, initialSize)}
}

func (e *engine) Set(_ context.Context, key, value string) error {
	e.data[key] = value
	return nil
}

func (e *engine) Get(_ context.Context, key string) (string, error) {
	v, ok := e.data[key]
	if !ok {
		return "", domain.ErrNotFound
	}

	return v, nil
}

func (e *engine) Delete(_ context.Context, key string) error {
	_, ok := e.data[key]
	if !ok {
		return domain.ErrNotFound
	}

	delete(e.data, key)
	return nil
}
