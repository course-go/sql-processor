package testlogger_test

import (
	"testing"

	"github.com/course-go/sql-processor/internal/test/testlogger"
)

func TestErrorLoggerWrites(t *testing.T) {
	t.Parallel()

	logger, writer := testlogger.NewTestErrorLogger()

	logger.Debug("doing")
	logger.Error("error")
	logger.Warn("logger")
	logger.Error("writes")
	logger.Info("testing")

	writer.AssertWrites(t, 2)
}
