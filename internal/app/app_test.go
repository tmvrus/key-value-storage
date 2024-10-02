package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tmvrus/key-value-storage/internal/config"
	"go.uber.org/mock/gomock"
)

func TestApp_Run(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	t.Run("do not start when address is invalid", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{}
		cfg.Network.Address = "invalid-host:1234"

		err := New(cfg, nil, log).Run(context.Background())
		require.Error(t, err)
		require.Contains(t, err.Error(), "no such host")

	})

	t.Run("stop when context is done", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{}
		cfg.Network.Address = findFreePort(t)

		cxt, cancel := context.WithCancel(context.Background())
		cancel()

		err := New(cfg, nil, log).Run(cxt)
		require.Error(t, err)
		require.True(t, errors.Is(err, context.Canceled))
	})

	t.Run("obey the limit", func(t *testing.T) {
		t.Parallel()

		cfg := &config.Config{}
		cfg.Network.Address = findFreePort(t)
		cfg.Network.MaxConnections = 1
		cfg.Network.IdleTimeout = time.Minute
		cfg.Network.MaxMessageSize = 1024

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)

		cxt, cancel := context.WithCancel(context.Background())

		stopped := make(chan struct{})
		go func() {
			err := New(cfg, storMock, log).Run(cxt)
			require.Error(t, err)
			require.True(t, errors.Is(err, context.Canceled))
			close(stopped)
		}()

		// wait for listening
		time.Sleep(1 * time.Second)

		conn1, err := net.Dial("tcp", cfg.Network.Address)
		checkConnectionOK(t, conn1, storMock)
		require.NoError(t, err)

		conn2, err := net.Dial("tcp", cfg.Network.Address)
		require.NoError(t, err)

		_, err = conn2.Read(make([]byte, 10))
		require.Error(t, err)
		require.True(t, errors.Is(err, io.EOF))

		require.NoError(t, conn1.Close())
		require.NoError(t, conn2.Close())

		conn1, err = net.Dial("tcp", cfg.Network.Address)
		require.NoError(t, err)
		checkConnectionOK(t, conn1, storMock)
		require.NoError(t, conn1.Close())

		cancel()
		<-stopped
	})
}

func checkConnectionOK(t *testing.T, c net.Conn, mock *Mockstorage) {
	t.Helper()

	mock.
		EXPECT().
		Set(gomock.Any(), "1", "1").
		Return(nil)

	cmd := []byte("SET 1 1\n")
	_, err := c.Write(cmd)
	require.NoError(t, err)

	response := make([]byte, 3)
	_, err = c.Read(response)
	require.NoError(t, err)
	require.Equal(t, "OK\n", string(response))
}

func findFreePort(t *testing.T) string {
	t.Helper()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)
	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	require.NoError(t, l.Close())

	return l.Addr().String()
}
