package wal

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/tmvrus/key-value-storage/internal/config"
	"github.com/tmvrus/key-value-storage/internal/domain"
)

type WAL struct {
	batchLock sync.Mutex
	batch     []record

	log     *slog.Logger
	batches chan []record

	walWriter writer

	flushingBatchSize    int
	flushingBatchTimeout time.Duration
}

type record struct {
	cmd domain.Command
	err chan error
}

func New(cfg config.Wal, w writer, log *slog.Logger) *WAL {
	return &WAL{
		log:                  log,
		walWriter:            w,
		batches:              make(chan []record, 1),
		flushingBatchTimeout: cfg.FlushingBatchTimeout,
		flushingBatchSize:    cfg.FlushingBatchSize,
	}
}

func (w *WAL) Add(_ context.Context, cmd domain.Command) error {
	rec := record{
		cmd: cmd,
		err: make(chan error),
	}

	withLock(&w.batchLock, func() {
		w.batch = append(w.batch, rec)
		if len(w.batch) >= w.flushingBatchSize {
			batch := w.batch
			w.batches <- batch
			w.batch = make([]record, 0, w.flushingBatchSize)
		}
	})

	return <-rec.err
}

func (w *WAL) Start(ctx context.Context) {
	ticker := time.NewTicker(w.flushingBatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.flush()
			w.log.DebugContext(ctx, "finish wal processing")
			return
		case <-ticker.C:
			w.flush()
		case batch := <-w.batches:
			w.walWriter.Write(batch)
			ticker.Reset(w.flushingBatchTimeout)
		}
	}
}

func (w *WAL) flush() {
	withLock(&w.batchLock, func() {
		batch := w.batch
		w.batch = make([]record, 0, w.flushingBatchSize)
		w.walWriter.Write(batch)
	})
}

func withLock(locker sync.Locker, doFunc func()) {
	locker.Lock()
	doFunc()
	locker.Unlock()
}
