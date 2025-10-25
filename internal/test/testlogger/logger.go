package testlogger

import (
	"io"
	"log/slog"
	"testing"
)

var _ io.Writer = &LoggerWriter{}

// LoggerWriter is a test logger writer implementation.
//
// Apart from satisfying the [io.Writer] interface it also exposes useful testing functions.
type LoggerWriter struct {
	writes int
}

func (lw *LoggerWriter) AssertWrites(t *testing.T, expectedWrites int) {
	t.Helper()

	if lw.writes != expectedWrites {
		t.Fatalf("unexpected logger writes: expected = %v, actual = %v", expectedWrites, lw.writes)
	}
}

// Write implements io.Writer.
func (lw *LoggerWriter) Write(p []byte) (n int, err error) {
	lw.writes++
	return len(p), nil
}

func NewTestErrorLogger() (logger *slog.Logger, writer *LoggerWriter) {
	writer = &LoggerWriter{}
	logger = slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelError}))
	return logger, writer
}
