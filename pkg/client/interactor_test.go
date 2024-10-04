package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConsoleInteractor(t *testing.T) {
	t.Parallel()

	t.Run("read ok", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)

		mock.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
			return copy(p, "ok\n"), nil
		})

		cmd, err := newConsoleInteractor(mock).ReadCommand()
		require.NoError(t, err)
		require.Equal(t, "ok", cmd)
	})

	t.Run("read error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)

		mock.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
			return 0, fmt.Errorf("ERROR")
		})

		_, err := newConsoleInteractor(mock).ReadCommand()
		require.Error(t, err)
	})

	t.Run("write ok", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)
		mock.EXPECT().Write([]byte("test")).Return(0, nil)

		err := newConsoleInteractor(mock).WriteResult("test")
		require.NoError(t, err)
	})

	t.Run("write error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mock := NewMockreaderWriter(ctrl)
		mock.EXPECT().Write([]byte("test")).Return(0, fmt.Errorf("ERROR"))

		err := newConsoleInteractor(mock).WriteResult("test")
		require.Error(t, err)
	})

}
