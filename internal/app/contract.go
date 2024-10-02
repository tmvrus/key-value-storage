//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -source=contract.go -destination=./contract_mock_test.go -package=app
package app

import (
	"net"

	stor "github.com/tmvrus/key-value-storage/internal/storage"
)

type storage interface {
	stor.Storage
}

type socket interface {
	net.Conn
}
