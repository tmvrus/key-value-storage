package storage

import (
	"context"

	"github.com/tmvrus/key-value-storage/internal/config"
	"github.com/tmvrus/key-value-storage/internal/storage/engine/inmemory"
)

type Storage interface {
	Set(cxt context.Context, key, value string) error
	Get(cxt context.Context, key string) (string, error)
	Delete(cxt context.Context, key string) error
}

func New(cfg *config.Config) Storage {
	switch cfg.Engine.Type {
	case config.EngineTypeInMemory:
		return inmemory.New()
	default:
		return inmemory.New()
	}
}
