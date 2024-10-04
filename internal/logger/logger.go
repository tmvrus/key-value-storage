package logger

import (
	"log/slog"
	"os"
	"time"
)

func New(fileName, level string) *slog.Logger {
	w := os.Stdout
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend|os.ModeExclusive)
	if err != nil {
		slog.
			Error("failed to open file for logging, use stdout", "error", err.Error(), "file_name", fileName)
	} else {
		go func() {
			t := time.NewTicker(time.Second)
			for range t.C {
				if err := f.Sync(); err != nil {
					slog.Error("failed to sync log file", "file_name", fileName, "error", err.Error())
				}
			}
		}()

		w = f
	}

	m := map[string]slog.Level{
		"info":  slog.LevelInfo,
		"debug": slog.LevelDebug,
		"error": slog.LevelError,
	}
	l, ok := m[level]
	if !ok {
		slog.Error("failed to find log level, use debug", "level", level)
		l = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: l}))
}
