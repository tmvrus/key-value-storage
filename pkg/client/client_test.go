package client

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_Client(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("loop stopped when got io.EOF error from readerWriter", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		interactorMock := NewMockreaderWriter(ctrl)

		hello := []byte("Waiting for command\n")
		interactorMock.EXPECT().Write(byteMatcher{t: t, want: hello}).Return(0, nil)
		interactorMock.EXPECT().Read(gomock.Any()).Return(0, io.EOF)

		client := NewClient(nil, log)
		client.StartInteractionLoop(interactorMock, interactorMock)
	})

	t.Run("read command, send and write result successfully", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		interactorMock := NewMockreaderWriter(ctrl)
		socketMock := NewMockreaderWriter(ctrl)

		hello := []byte("Waiting for command\n")
		interactorMock.EXPECT().Write(byteMatcher{t: t, want: hello}).Return(0, nil)

		cmd := []byte("GET KEY\n")
		interactorMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, cmd), io.EOF
			}).Times(1)

		matcher := byteMatcher{t: t, want: cmd}
		socketMock.EXPECT().Write(matcher).Return(len(cmd), nil)

		response := []byte("VALUE")

		socketMock.
			EXPECT().
			Read(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return copy(p, response), nil
			})

		interactorMock.EXPECT().Write(byteMatcher{t: t, want: response}).Return(len(response), nil)
		interactorMock.EXPECT().Write(byteMatcher{t: t, want: []byte("\n")}).Return(1, nil)

		client := NewClient(socketMock, log)
		client.StartInteractionLoop(interactorMock, interactorMock)
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
