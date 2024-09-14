package app

import (
	"context"
)

type storage interface {
	Set(cxt context.Context, key, value string) error
	Get(cxt context.Context, key string) (string, error)
	Delete(cxt context.Context, key string) error
}
