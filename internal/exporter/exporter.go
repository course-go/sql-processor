package exporter

import "github.com/course-go/sql-processor/internal/sql"

// Exporter represents [sql.Statement] exporter.
type Exporter interface {
	Export(statement sql.Statement) (err error)
	ExportBatch(statements []sql.Statement) (err error)
}
