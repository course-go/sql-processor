package processor_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/course-go/sql-processor/internal/processor"
	"github.com/course-go/sql-processor/internal/sql"
	"github.com/course-go/sql-processor/internal/test/testlogger"
)

const filePermissions = 0o755

func TestRun(t *testing.T) { //nolint: cyclop, gocognit, maintidx
	t.Parallel()

	t.Run("NonexistentFile", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, loggerWriter := testlogger.NewTestErrorLogger()
			fileCh := make(chan sql.File, 1)
			statementCh := make(chan sql.Statement, 1)

			p := processor.New(logger, fileCh, statementCh)

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go p.Run(ctx)

			fileCh <- sql.File{
				Path: filepath.Join(t.TempDir(), "nonexistent.sql"),
				Type: sql.SQLite,
			}

			// Wait for processor to process the file.
			synctest.Wait()

			select {
			case statement := <-statementCh:
				t.Errorf("unexpected statement in channel: expected nothing, got = %v", statement)
			default:
			}

			loggerWriter.AssertWrites(t, 1)
		})
	})

	t.Run("EmptyFile", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, loggerWriter := testlogger.NewTestErrorLogger()
			fileCh := make(chan sql.File, 1)
			statementCh := make(chan sql.Statement, 1)

			p := processor.New(logger, fileCh, statementCh)

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go p.Run(ctx)

			path := filepath.Join(t.TempDir(), "file.sql")
			f, err := os.Create(path)
			if err != nil {
				t.Errorf("failed creating test file: %v", err)
			}

			defer func() {
				_ = f.Close()
			}()

			fileCh <- sql.File{
				Path: path,
				Type: sql.SQLite,
			}

			// Wait for processor to process the file.
			synctest.Wait()

			select {
			case statement := <-statementCh:
				t.Errorf("unexpected statement in channel: expected nothing, got = %v", statement)
			default:
			}

			loggerWriter.AssertWrites(t, 0)
		})
	})

	t.Run("SimpleSQLFile", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, loggerWriter := testlogger.NewTestErrorLogger()
			fileCh := make(chan sql.File, 1)
			statementCh := make(chan sql.Statement, 1)

			p := processor.New(logger, fileCh, statementCh)

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go p.Run(ctx)

			path := copyFile(t, t.TempDir(), filepath.Join("testdata", "test.sql"))
			file := sql.File{
				Path: path,
				Type: sql.SQLite,
			}
			fileCh <- file

			var statements []sql.Statement
			go func() {
				for {
					select {
					case statement := <-statementCh:
						statements = append(statements, statement)
					case <-ctx.Done():
						return
					}
				}
			}()

			// Wait for processor to process the file.
			synctest.Wait()

			expectedStatements := []sql.Statement{
				{
					File:    file,
					Content: "SELECT * FROM users",
					LineNum: 1,
				},
				{
					File:    file,
					Content: "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
					LineNum: 2,
				},
				{
					File:    file,
					Content: "UPDATE users SET name = 'Jane' WHERE id = 1",
					LineNum: 3,
				},
				{
					File:    file,
					Content: "DELETE FROM users WHERE id = 2",
					LineNum: 4,
				},
			}

			loggerWriter.AssertWrites(t, 0)

			if len(statements) != len(expectedStatements) {
				t.Fatalf(
					"unexpected statements count: expected = %v, got = %v",
					len(expectedStatements),
					len(statements),
				)
			}

			for i := range statements {
				if statements[i] != expectedStatements[i] {
					t.Errorf("statement does not match: expected = %v, got = %v", expectedStatements[i], statements[i])
				}
			}
		})
	})

	t.Run("SimpleSQLFileWithComments", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, loggerWriter := testlogger.NewTestErrorLogger()
			fileCh := make(chan sql.File, 1)
			statementCh := make(chan sql.Statement, 1)

			p := processor.New(logger, fileCh, statementCh)

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go p.Run(ctx)

			path := copyFile(t, t.TempDir(), filepath.Join("testdata", "test-comments.sql"))
			file := sql.File{
				Path: path,
				Type: sql.SQLite,
			}
			fileCh <- file

			var statements []sql.Statement
			go func() {
				for {
					select {
					case statement := <-statementCh:
						statements = append(statements, statement)
					case <-ctx.Done():
						return
					}
				}
			}()

			// Wait for processor to process the file.
			synctest.Wait()

			expectedStatements := []sql.Statement{
				{
					File:    file,
					Content: "SELECT * FROM users",
					LineNum: 2,
				},
				{
					File:    file,
					Content: "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
					LineNum: 4,
				},
				{
					File:    file,
					Content: "UPDATE users SET name = 'Jane' WHERE id = 1",
					LineNum: 7,
				},
				{
					File:    file,
					Content: "DELETE FROM users WHERE id = 2",
					LineNum: 9,
				},
			}

			loggerWriter.AssertWrites(t, 0)

			if len(statements) != len(expectedStatements) {
				t.Fatalf(
					"unexpected statements count: expected = %v, got = %v",
					len(expectedStatements),
					len(statements),
				)
			}

			for i := range statements {
				if statements[i] != expectedStatements[i] {
					t.Errorf("statement does not match: expected = %v, got = %v", expectedStatements[i], statements[i])
				}
			}
		})
	})

	t.Run("MultilineQLFileWithComments", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, loggerWriter := testlogger.NewTestErrorLogger()
			fileCh := make(chan sql.File, 1)
			statementCh := make(chan sql.Statement, 1)

			p := processor.New(logger, fileCh, statementCh)

			ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
			defer cancel()

			go p.Run(ctx)

			path := copyFile(t, t.TempDir(), filepath.Join("testdata", "test-multiline.sql"))
			file := sql.File{
				Path: path,
				Type: sql.SQLite,
			}
			fileCh <- file

			var statements []sql.Statement
			go func() {
				for {
					select {
					case statement := <-statementCh:
						statements = append(statements, statement)
					case <-ctx.Done():
						return
					}
				}
			}()

			// Wait for processor to process the file.
			synctest.Wait()

			expectedStatements := []sql.Statement{
				{
					File:    file,
					Content: "SELECT * FROM users",
					LineNum: 2,
				},
				{
					File:    file,
					Content: "INSERT INTO users (name, email)\nVALUES ('John', 'john@example.com')",
					LineNum: 4,
				},
				{
					File:    file,
					Content: "UPDATE users\nSET name = 'Jane'\nWHERE id = 1",
					LineNum: 9,
				},
				{
					File:    file,
					Content: "UPDATE users\nSET name = 'Bob'\nWHERE id = 4",
					LineNum: 14,
				},

				{
					File:    file,
					Content: "DELETE FROM users WHERE id = 2",
					LineNum: 19,
				},
			}

			loggerWriter.AssertWrites(t, 0)

			if len(statements) != len(expectedStatements) {
				t.Fatalf(
					"unexpected statements count: expected = %v, got = %v",
					len(expectedStatements),
					len(statements),
				)
			}

			for i := range statements {
				if statements[i] != expectedStatements[i] {
					t.Errorf("statement does not match: expected = %v, got = %v", expectedStatements[i], statements[i])
				}
			}
		})
	})
}

func copyFile(t *testing.T, directory string, file string) string {
	t.Helper()

	bytes, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed reading testdata file %v: %v", file, err)
	}

	path := filepath.Join(directory, filepath.Base(file))
	err = os.WriteFile(path, bytes, filePermissions)
	if err != nil {
		t.Fatalf("failed writing to temp file: %v", err)
	}

	// Simulate delay.
	time.Sleep(100 * time.Millisecond)

	return path
}
