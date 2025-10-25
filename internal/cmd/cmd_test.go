package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/course-go/sql-processor/internal/cmd"
	"github.com/course-go/sql-processor/internal/exporter"
	"github.com/course-go/sql-processor/internal/sql"
	"github.com/course-go/sql-processor/internal/test/testexporter"
)

const filePermissions = 0o755

func TestRun(t *testing.T) { //nolint: cyclop, gocognit
	t.Parallel()

	testdataPostgresFilesDirectory := filepath.Join("testdata", "postgres")
	postgresFiles := []string{
		filepath.Join(testdataPostgresFilesDirectory, "test-select.sql"),
		filepath.Join(testdataPostgresFilesDirectory, "test-update.sql"),
		filepath.Join(testdataPostgresFilesDirectory, "test-joins.sql"),
	}
	postgresStatementCount := 11

	testdataMySQLFilesDirectory := filepath.Join("testdata", "mysql")
	mysqlFiles := []string{
		filepath.Join(testdataMySQLFilesDirectory, "test-select.sql"),
		filepath.Join(testdataMySQLFilesDirectory, "test-update.sql"),
		filepath.Join(testdataMySQLFilesDirectory, "test-joins.sql"),
		filepath.Join(testdataMySQLFilesDirectory, "test-create.sql"),
	}
	mysqlStatementCount := 11

	t.Run("NoArguments", func(t *testing.T) {
		t.Parallel()

		args := []string{"sql-processor"}
		err := cmd.Run(t.Context(), args, nil)
		if err == nil {
			t.Fatalf("expected error for empty arguments")
		}
	})

	t.Run("SingleDirectoryDirective", func(t *testing.T) {
		t.Parallel()

		directory := filepath.Join(t.TempDir(), "postgres")
		createDirectory(t, directory)

		directive := directory + ":" + string(sql.PostgresType)
		args := []string{"sql-processor", directive}

		e := testexporter.New()
		exporters := []exporter.Exporter{e}

		errCh := make(chan error, 1)
		go func() {
			err := cmd.Run(t.Context(), args, exporters)
			errCh <- err
		}()

		// Unfortunate.
		// There is no reasonable way to handle this without polluting the code.
		time.Sleep(1 * time.Second)

		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("run returned unexpected error: %v", err)
			}
		default:
			// Run should block the goroutine so nothing should really be returned.
		}

		copyFiles(t, directory, postgresFiles)

		ticker := time.NewTicker(10 * time.Millisecond)
		deadline := time.NewTimer(3 * time.Second)

		for {
			select {
			case <-ticker.C:
				if len(e.Statements()) == postgresStatementCount {
					return
				}

				if len(e.Statements()) > postgresStatementCount {
					t.Fatalf(
						"received more statements than expected: expected = %v, got = %v",
						postgresStatementCount,
						len(e.Statements()),
					)
				}

				t.Logf(
					"checking statement count: expected = %v, current = %v",
					postgresStatementCount,
					len(e.Statements()),
				)

			case <-deadline.C:
				t.Fatalf("failed processing data in time")
			}
		}
	})

	t.Run("MultipleDirectoryDirectives", func(t *testing.T) {
		t.Parallel()

		postgresDirectory := filepath.Join(t.TempDir(), "postgres")
		createDirectory(t, postgresDirectory)

		mysqlDirectory := filepath.Join(t.TempDir(), "mysql")
		createDirectory(t, mysqlDirectory)

		postgresDirective := postgresDirectory + ":" + string(sql.PostgresType)
		mysqlDirective := mysqlDirectory + ":" + string(sql.MySQL)
		args := []string{"sql-processor", postgresDirective, mysqlDirective}

		e := testexporter.New()
		exporters := []exporter.Exporter{e}

		errCh := make(chan error, 1)
		go func() {
			err := cmd.Run(t.Context(), args, exporters)
			errCh <- err
		}()

		// Unfortunate.
		// There is no reasonable way to handle this without polluting the code.
		time.Sleep(1 * time.Second)

		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("run returned unexpected error: %v", err)
			}
		default:
			// Run should block the goroutine so nothing should really be returned.
		}

		copyFiles(t, postgresDirectory, postgresFiles)
		copyFiles(t, mysqlDirectory, mysqlFiles)

		totalStatementCount := postgresStatementCount + mysqlStatementCount
		ticker := time.NewTicker(10 * time.Millisecond)
		deadline := time.NewTimer(3 * time.Second)

		for {
			select {
			case <-ticker.C:
				if len(e.Statements()) == totalStatementCount {
					return
				}

				if len(e.Statements()) > totalStatementCount {
					t.Fatalf(
						"received more statements than expected: expected = %v, got = %v",
						totalStatementCount,
						len(e.Statements()),
					)
				}

				t.Logf(
					"checking statement count: expected = %v, current = %v",
					totalStatementCount,
					len(e.Statements()),
				)

			case <-deadline.C:
				t.Fatalf("failed processing data in time")
			}
		}
	})
}

func createDirectory(t *testing.T, path string) {
	t.Helper()

	err := os.MkdirAll(path, filePermissions)
	if err != nil {
		t.Fatalf("failed creating %v test directory: %v", path, err)
	}
}

func copyFiles(t *testing.T, directory string, files []string) {
	t.Helper()

	for _, file := range files {
		bytes, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("failed reading testdata file %v: %v", file, err)
		}

		err = os.WriteFile(filepath.Join(directory, filepath.Base(file)), bytes, filePermissions)
		if err != nil {
			t.Fatalf("failed writing to temp file: %v", err)
		}

		// Simulate delay.
		time.Sleep(100 * time.Millisecond)
	}
}
