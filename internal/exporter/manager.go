package exporter

import (
	"context"
	"log/slog"

	"github.com/course-go/sql-processor/internal/sql"
)

// Manager manages [Exporter]s.
// It listens for processed [sql.Statement]s and passes them down to all exporters for exporting.
type Manager struct {
	logger      *slog.Logger
	statementCh <-chan sql.Statement
	exporters   []Exporter
}

func NewManager(logger *slog.Logger, statementCh <-chan sql.Statement, exporters []Exporter) Manager {
	return Manager{
		logger:      logger.With("component", "exporter-manager"),
		statementCh: statementCh,
		exporters:   exporters,
	}
}

// Run runs the [Manager].
func (m *Manager) Run(ctx context.Context) {
	// TODO: Implement.
}
