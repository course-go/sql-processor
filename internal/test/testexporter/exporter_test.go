package testexporter_test

import (
	"testing"

	"github.com/course-go/sql-processor/internal/sql"
	"github.com/course-go/sql-processor/internal/test/testexporter"
)

func TestExporter(t *testing.T) {
	t.Parallel()

	e := testexporter.New()

	statements := []sql.Statement{
		{
			File: sql.File{
				Path: "test.sql",
				Type: sql.PostgresType,
			},
			Content: "SELECT * FROM USERS",
			LineNum: 1,
		},
		{
			File: sql.File{
				Path: "test.sql",
				Type: sql.PostgresType,
			},
			Content: "DROP USERS",
			LineNum: 2,
		},
	}

	err := e.ExportBatch(statements)
	if err != nil {
		t.Fatalf("failed exporting batch: %v", err)
	}

	if len(statements) != len(e.Statements()) {
		t.Fatalf("statement lengths do not match: expected = %v, got = %v", len(statements), len(e.Statements()))
	}

	for i := range statements {
		if statements[i] != e.Statements()[i] {
			t.Fatalf("statements do not match: expected = %v, got = %v", statements[i], e.Statements()[i])
		}
	}
}
