package exporter_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/course-go/sql-processor/internal/exporter"
	"github.com/course-go/sql-processor/internal/sql"
	"github.com/course-go/sql-processor/internal/test/testexporter"
	"github.com/course-go/sql-processor/internal/test/testlogger"
)

func TestManager(t *testing.T) {
	t.Parallel()

	file1 := sql.File{
		Path: "test1.sql",
		Type: "mysql",
	}
	file2 := sql.File{
		Path: "test2.sql",
		Type: "mysql",
	}
	statements := []sql.Statement{
		{Content: "SELECT * FROM users", LineNum: 1, File: file1},
		{Content: "INSERT INTO users VALUES (1, 'John')", LineNum: 2, File: file1},
		{Content: "UPDATE users SET name = 'Jane'", LineNum: 1, File: file2},
	}

	t.Run("SingleExporter", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			mock := testexporter.New()
			statementCh := make(chan sql.Statement, 1)
			logger, _ := testlogger.NewTestErrorLogger()
			p := exporter.NewManager(logger, statementCh, []exporter.Exporter{mock})

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go p.Run(ctx)

			// Send test statements to exporter manager.
			for _, statement := range statements {
				statementCh <- statement
			}

			// Wait for exporter to receive data or timeout.
			synctest.Wait()

			// Check that exporter statements match.
			for i := range statements {
				if mock.Statements()[i] != statements[i] {
					t.Fatalf(
						"exporter data does not match: expected = %v, got = %v",
						statements[i],
						mock.Statements()[i],
					)
				}
			}
		})
	})

	t.Run("MultipleExporters", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			mocks := []*testexporter.Exporter{
				testexporter.New(),
				testexporter.New(),
				testexporter.New(),
			}

			statementCh := make(chan sql.Statement, 1)

			var exporters []exporter.Exporter
			for _, mock := range mocks {
				exporters = append(exporters, mock)
			}

			logger, _ := testlogger.NewTestErrorLogger()
			m := exporter.NewManager(logger, statementCh, exporters)

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go m.Run(ctx)

			// Send test statements to exporter manager.
			for _, statement := range statements {
				statementCh <- statement
			}

			// Wait for exporters to receive data or timeout.
			synctest.Wait()

			// Check that exporter statements match.
			for _, mock := range mocks {
				for i := range statements {
					if mock.Statements()[i] != statements[i] {
						t.Fatalf(
							"exporter data does not match: expected = %v, got = %v",
							statements[i],
							mock.Statements()[i],
						)
					}
				}
			}
		})
	})
}
