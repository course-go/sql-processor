package stdout

import (
	"errors"
	"fmt"

	"github.com/course-go/sql-processor/internal/exporter"
	"github.com/course-go/sql-processor/internal/sql"
)

var _ exporter.Exporter = &Exporter{}

// Exporter implements [exporter.Exporter] and exports given [sql.Statement]s to stdout.
type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

// Export implements exporter.Exporter.
func (e *Exporter) Export(statement sql.Statement) (err error) {
	fmt.Printf("%s:%d [%s] [%s]\n", statement.File.Path, statement.LineNum, statement.File.Type, statement.Content)
	return nil
}

// ExportBatch implements exporter.Exporter.
func (e *Exporter) ExportBatch(statements []sql.Statement) (err error) {
	for _, statement := range statements {
		err = errors.Join(err, e.Export(statement))
	}

	return err
}
