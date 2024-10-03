package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestApp_Session(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := handlerConfig{
		timeout:    time.Minute,
		bufferSize: 1024,
	}

	t.Run("handle canceled context", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		newHandler(log, nil, nil, cfg).startHandling(ctx)
	})

	t.Run("handle SetReadDeadline error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(fmt.Errorf("ERROR"))

		newHandler(log, nil, socketMock, cfg).startHandling(ctx)
	})

	t.Run("handle deadline/idle on read", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(nil)

		socketMock.
			EXPECT().
			Read(gomock.Any()).
			Return(0, os.ErrDeadlineExceeded)

		newHandler(log, nil, socketMock, cfg).startHandling(ctx)
	})

	t.Run("able to read command, execute and write response", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(nil).Times(2)

		cmd := []byte("DELETE KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		storMock.EXPECT().Delete(ctx, "KEY").Return(nil)

		socketMock.EXPECT().SetWriteDeadline(inFuture{t}).Return(nil)
		socketMock.EXPECT().Write(byteMatcher{t: t, want: []byte("OK\n")}).Return(0, nil)

		newHandler(log, storMock, socketMock, cfg).startHandling(ctx)
	})

	t.Run("able to fail conn read fail", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(nil)

		socketMock.
			EXPECT().
			Read(gomock.Any()).
			Return(0, io.EOF).
			Times(1)

		newHandler(log, storMock, socketMock, cfg).startHandling(ctx)
	})

	t.Run("able to handle invalid command", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(nil).Times(2)
		cmd := []byte("INVALID KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		socketMock.EXPECT().SetWriteDeadline(inFuture{t}).Return(nil)
		socketMock.EXPECT().Write(byteMatcher{t: t, want: []byte("ERROR: unsupported operation\n")}).Return(0, nil)

		newHandler(log, storMock, socketMock, cfg).startHandling(ctx)
	})

	t.Run("handle SetWriteDeadline error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(nil).Times(2)
		cmd := []byte("GET KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		storMock.EXPECT().Get(ctx, "KEY").Return("", fmt.Errorf("STORAGE"))

		socketMock.EXPECT().SetWriteDeadline(inFuture{t}).Return(fmt.Errorf("ERROR"))

		newHandler(log, storMock, socketMock, cfg).startHandling(ctx)
	})

	t.Run("able to handle storage error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		storMock := NewMockstorage(ctrl)
		socketMock := NewMocksocket(ctrl)
		ctx := context.Background()

		socketMock.EXPECT().SetReadDeadline(inFuture{t}).Return(nil).Times(2)
		cmd := []byte("GET KEY\n")
		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		storMock.EXPECT().Get(ctx, "KEY").Return("", fmt.Errorf("STORAGE"))

		socketMock.EXPECT().SetWriteDeadline(inFuture{t}).Return(nil)
		socketMock.EXPECT().Write(byteMatcher{t: t, want: []byte("ERROR: STORAGE\n")}).Return(0, nil)

		newHandler(log, storMock, socketMock, cfg).startHandling(ctx)
	})
}

type inFuture struct {
	t *testing.T
}

func (m inFuture) Matches(x any) bool {
	got, ok := x.(time.Time)
	if !ok {
		m.t.Errorf("expected time.Time got %T", x)
		return false
	}

	if n := time.Now(); !got.After(n) {
		m.t.Errorf("expected %s after now %s", got.String(), n)
		return false
	}

	return true
}

func (m inFuture) String() string {
	return "expected passed time be in future"
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
