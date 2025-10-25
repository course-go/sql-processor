package processor

import (
	"context"
	"log/slog"

	"github.com/course-go/sql-processor/internal/sql"
)

// Processor is a component that receives given [sql.File] and processes them to [sql.Statement]s.
// It reads the given files and parses the statements from them.
type Processor struct {
	logger      *slog.Logger
	fileCh      <-chan sql.File
	statementCh chan<- sql.Statement
}

func New(logger *slog.Logger, fileCh <-chan sql.File, statementCh chan<- sql.Statement) Processor {
	return Processor{
		logger:      logger.With("component", "processor"),
		fileCh:      fileCh,
		statementCh: statementCh,
	}
}

// Run runs the [Processor].
func (p *Processor) Run(ctx context.Context) {
	// TODO: Implement.
}
