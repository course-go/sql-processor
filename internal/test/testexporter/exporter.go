package testexporter

import (
	"errors"
	"sync"

	"github.com/course-go/sql-processor/internal/sql"
)

// Exporter is a test exporter implementation.
//
// Apart from satisfying the [exporter.Exporter] interface it also stores
// passed statements and gives the ability to access it.
type Exporter struct {
	mu         sync.Mutex
	statements []sql.Statement
}

func New() *Exporter {
	return &Exporter{}
}

// Export implements Exporter.
func (e *Exporter) Export(statement sql.Statement) (err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.statements = append(e.statements, statement)
	return nil
}

// ExportBatch implements Exporter.
func (e *Exporter) ExportBatch(statements []sql.Statement) (err error) {
	for _, statement := range statements {
		err = errors.Join(err, e.Export(statement))
	}

	return err
}

func (e *Exporter) Statements() []sql.Statement {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.statements
}
