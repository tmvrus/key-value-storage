package app

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestApp_Run(t *testing.T) {
	t.Parallel()

	t.Run("fail when parse fail", func(t *testing.T) {
		t.Parallel()

		in := bytes.NewReader([]byte("Invalid command"))
		out := bytes.NewBuffer(nil)
		New(
			in,
			out,
			nil,
		).Run(context.Background())

		require.Equal(t, "ERROR: invalid command Invalid for size 2\n", out.String())
	})

	t.Run("fail when cmd fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		in := bytes.NewReader([]byte("SET key value"))
		out := bytes.NewBuffer(nil)

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		mock := NewMockstorage(ctrl)

		mock.EXPECT().Set(ctx, "key", "value").Return(fmt.Errorf("STORAGE ERROR"))
		New(
			in,
			out,
			mock,
		).Run(ctx)

		require.Contains(t, out.String(), "STORAGE ERROR")
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		in := bytes.NewReader([]byte("GET key"))
		out := bytes.NewBuffer(nil)

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		mock := NewMockstorage(ctrl)

		mock.EXPECT().Get(ctx, "key").Return("value", nil)
		New(
			in,
			out,
			mock,
		).Run(ctx)

		require.Equal(t, "value\n", out.String())
	})
}
