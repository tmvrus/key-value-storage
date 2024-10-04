package client

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Sender(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)

		mock.EXPECT().Write([]byte("test")).Return(0, nil)
		mock.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
			return copy(p, []byte("ok")), nil
		})

		res, err := newSender(mock).Send(ctx, "test")
		require.NoError(t, err)
		require.Equal(t, "ok", res)
	})

	t.Run("handle write error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)

		mock.EXPECT().Write([]byte("test")).Return(0, fmt.Errorf("ERROR"))

		_, err := newSender(mock).Send(ctx, "test")
		require.Error(t, err)
	})

	t.Run("handle read error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)

		mock.EXPECT().Write([]byte("test")).Return(0, nil)
		mock.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
			return 0, fmt.Errorf("ERROR")
		})

		_, err := newSender(mock).Send(ctx, "test")
		require.Error(t, err)
	})

	t.Run("handle context canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)

		mock.EXPECT().Write([]byte("test")).Return(0, nil).AnyTimes()

		_, err := newSender(mock).Send(ctx, "test")
		require.True(t, errors.Is(err, context.Canceled))
	})

}
