package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestApp_Session(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("able to read command, execute and write response", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		cmd := []byte("DELETE KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		storMock.EXPECT().Delete(ctx, "KEY").Return(nil)

		socketMock.EXPECT().Write(byteMatcher{t: t, want: []byte("OK\n")}).Return(0, nil)

		newSession(log, storMock, socketMock).start(ctx)
	})

	t.Run("able to fail socket read fail", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.
			EXPECT().
			Read(gomock.Any()).
			Return(0, io.EOF).
			Times(1)

		newSession(log, storMock, socketMock).start(ctx)
	})

	t.Run("able to handle invalid command", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		cmd := []byte("INVALID KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		socketMock.EXPECT().Write(byteMatcher{t: t, want: []byte("ERROR: unsupported operation\n")}).Return(0, nil)

		newSession(log, storMock, socketMock).start(ctx)
	})

	t.Run("able to handle storage error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		cmd := []byte("GET KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		storMock.EXPECT().Get(ctx, "KEY").Return("", fmt.Errorf("STORAGE"))

		socketMock.EXPECT().Write(byteMatcher{t: t, want: []byte("ERROR: STORAGE\n")}).Return(0, nil)

		newSession(log, storMock, socketMock).start(ctx)
	})

}

type byteMatcher struct {
	t    *testing.T
	want []byte
}

func (m byteMatcher) Matches(x any) bool {
	got, ok := x.([]byte)
	if !ok {
		m.t.Errorf("expected []byte got %T", x)
		return false
	}

	gotS := string(got)
	wantS := string(m.want)
	if gotS != wantS {
		m.t.Errorf("expected %q, but got %q", wantS, gotS)
		return false
	}

	return true
}

func (m byteMatcher) String() string {
	return ""
}
