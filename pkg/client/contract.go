//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=client
package client

import "io"

type readerWriter interface {
	io.Writer
	io.Reader
}
