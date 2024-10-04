package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
	"testing"

	"go.uber.org/mock/gomock"
)

func setupClient(t *testing.T) (*Client, *Mockinteractor, *Mocksender) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	senderMock := NewMocksender(ctrl)
	interactorMock := NewMockinteractor(ctrl)
	return &Client{
		sender:     senderMock,
		interactor: interactorMock,
		log:        slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}, interactorMock, senderMock

}

func Test_Client(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("loop stopped when got io.EOF error from readerWriter", func(t *testing.T) {
		t.Parallel()

		client, im, _ := setupClient(t)

		im.EXPECT().ReadCommand().Return("", fmt.Errorf("not critical error"))
		im.EXPECT().ReadCommand().Return("", fmt.Errorf("critical: %w", io.EOF))

		client.Start(ctx)
	})

	t.Run("read command, send and write result successfully", func(t *testing.T) {
		t.Parallel()

		client, im, sm := setupClient(t)

		im.EXPECT().ReadCommand().Return("COMMAND", nil)
		sm.EXPECT().Send(ctx, "COMMAND").Return("RESULT", nil)
		im.EXPECT().WriteResult("RESULT").Return(nil)

		im.EXPECT().ReadCommand().Return("", io.EOF)

		client.Start(ctx)
	})

	t.Run("fail when send fails", func(t *testing.T) {
		t.Parallel()

		client, im, sm := setupClient(t)

		im.EXPECT().ReadCommand().Return("COMMAND", nil)
		sm.EXPECT().Send(ctx, "COMMAND").Return("", syscall.EPIPE)

		client.Start(ctx)
	})

	t.Run("fail when write result fails", func(t *testing.T) {
		t.Parallel()

		client, im, sm := setupClient(t)

		im.EXPECT().ReadCommand().Return("COMMAND", nil)
		sm.EXPECT().Send(ctx, "COMMAND").Return("RESULT", nil)
		im.EXPECT().WriteResult("RESULT").Return(syscall.ECONNRESET)

		client.Start(ctx)
	})
}
