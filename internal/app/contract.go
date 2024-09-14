//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=app
package app

import (
	"context"
)

type storage interface {
	Set(cxt context.Context, key, value string) error
	Get(cxt context.Context, key string) (string, error)
	Delete(cxt context.Context, key string) error
}
