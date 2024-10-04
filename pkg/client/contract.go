//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=client
package client

import (
	"context"
	"io"
)

type readerWriter interface {
	io.Writer
	io.Reader
}

type interactor interface {
	ReadCommand() (string, error)
	WriteResult(s string) error
}

type sender interface {
	Send(context.Context, string) (string, error)
}
