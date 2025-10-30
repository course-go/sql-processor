package observer_test

import (
	"cmp"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/course-go/sql-processor/internal/observer"
	"github.com/course-go/sql-processor/internal/sql"
	"github.com/course-go/sql-processor/internal/test/testlogger"
)

func TestObserver(t *testing.T) { //nolint: cyclop, gocognit, maintidx
	t.Parallel()

	postgresFiles := []sql.File{
		{
			Path: "test1.sql",
			Type: sql.PostgresType,
		},
		{
			Path: "test2.sql",
			Type: sql.PostgresType,
		},
	}
	mysqlFiles := []sql.File{
		{
			Path: "test3.sql",
			Type: sql.MySQL,
		},
		{
			Path: "test4.sql",
			Type: sql.MySQL,
		},
		{
			Path: "test5.sql",
			Type: sql.MySQL,
		},
	}

	t.Run("EmptyDirectives", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, _ := testlogger.NewTestErrorLogger()
			_, err := observer.New(logger, []string{}, make(chan sql.File))
			if !errors.Is(err, observer.ErrNoDirectoryDirectivesProvided) {
				t.Fatalf("expected no directives provided error: got = %v", err)
			}
		})
	})

	t.Run("InvalidDirective", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, _ := testlogger.NewTestErrorLogger()
			_, err := observer.New(logger, []string{"/path/to/dir"}, make(chan sql.File))
			if !errors.Is(err, observer.ErrInvalidDirectoryDirective) {
				t.Fatalf("expected invalid directive error: got = %v", err)
			}
		})
	})

	t.Run("UnknownDirectiveType", func(t *testing.T) {
		t.Parallel()

		synctest.Test(t, func(t *testing.T) {
			logger, _ := testlogger.NewTestErrorLogger()
			_, err := observer.New(logger, []string{"/path/to/dir:unknown-type"}, make(chan sql.File))
			if !errors.Is(err, sql.ErrUnknownType) {
				t.Fatalf("expected unknown type error: got = %v", err)
			}
		})
	})

	t.Run("SingleDirectory", func(t *testing.T) {
		t.Parallel()

		fileCh := make(chan sql.File, len(postgresFiles))
		logger, _ := testlogger.NewTestErrorLogger()

		directory := filepath.Join(t.TempDir(), "postgres")
		err := os.Mkdir(directory, 0o700)
		if err != nil {
			t.Fatalf("failed to create file directory: %v", err)
		}

		directive := directory + ":" + string(sql.PostgresType)
		o, err := observer.New(logger, []string{directive}, fileCh)
		if err != nil {
			t.Fatalf("failed to create observer: %v", err)
		}

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		go o.Run(ctx)

		defer func() {
			err := o.Close()
			if err != nil {
				t.Fatalf("failed to close observer: %v", err)
			}
		}()

		var files []sql.File
		fileReaderDone := make(chan bool, 1)
		go func() {
			defer func() {
				fileReaderDone <- true
			}()

			for {
				select {
				case file := <-fileCh:
					files = append(files, file)
					if len(files) == len(postgresFiles) {
						return
					}
				case <-time.After(5 * time.Second): // Timeout unless the reader finishes.
					return
				}
			}
		}()

		createFiles(t, directory, postgresFiles)

		<-fileReaderDone

		if len(files) != len(postgresFiles) {
			t.Errorf("observed file count does not match: expected = %v, got = %v", len(postgresFiles), len(files))
		}

		slices.SortFunc(files, func(a, b sql.File) int {
			return cmp.Compare(a.Path, b.Path)
		})

		postgresFiles := prefixPaths(directory, postgresFiles)

		// Check that exporter statements match.
		for i := range files {
			if files[i] != postgresFiles[i] {
				t.Errorf("received file does not match: expected = %v, got = %v", postgresFiles[i], files[i])
			}
		}
	})

	t.Run("MultipleDirectories", func(t *testing.T) {
		t.Parallel()

		fileCh := make(chan sql.File, len(postgresFiles))
		logger, _ := testlogger.NewTestErrorLogger()

		postgresDirectory := filepath.Join(t.TempDir(), "postgres")
		err := os.Mkdir(postgresDirectory, 0o700)
		if err != nil {
			t.Fatalf("failed to create file directory: %v", err)
		}

		mysqlDirectory := filepath.Join(t.TempDir(), "mysql")
		err = os.Mkdir(mysqlDirectory, 0o700)
		if err != nil {
			t.Fatalf("failed to create file directory: %v", err)
		}

		postgresDirective := postgresDirectory + ":" + string(sql.PostgresType)
		mysqlDirective := mysqlDirectory + ":" + string(sql.MySQL)
		o, err := observer.New(logger, []string{postgresDirective, mysqlDirective}, fileCh)
		if err != nil {
			t.Fatalf("failed to create observer: %v", err)
		}

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		go o.Run(ctx)

		defer func() {
			err := o.Close()
			if err != nil {
				t.Fatalf("failed to close observer: %v", err)
			}
		}()

		var files []sql.File
		fileReaderDone := make(chan bool, 1)
		go func() {
			defer func() {
				fileReaderDone <- true
			}()

			for {
				select {
				case file := <-fileCh:
					files = append(files, file)
					if len(files) == len(postgresFiles)+len(mysqlFiles) {
						return
					}
				case <-time.After(5 * time.Second): // Timeout unless the reader finishes.
					return
				}
			}
		}()

		createFiles(t, postgresDirectory, postgresFiles)
		createFiles(t, mysqlDirectory, mysqlFiles)

		<-fileReaderDone

		postgresFiles := prefixPaths(postgresDirectory, postgresFiles)
		mysqlFiles := prefixPaths(mysqlDirectory, mysqlFiles)
		expectedFiles := append(postgresFiles, mysqlFiles...) //nolint: gocritic

		if len(files) != len(expectedFiles) {
			t.Errorf(
				"observed file count does not match: expected = %v, got = %v",
				len(expectedFiles),
				len(files),
			)
		}

		slices.SortFunc(files, func(a, b sql.File) int {
			return cmp.Compare(a.Path, b.Path)
		})

		// Check that exporter statements match.
		for i := range files {
			if files[i] != expectedFiles[i] {
				t.Errorf("received file does not match: expected = %v, got = %v", expectedFiles[i], files[i])
			}
		}
	})
}

func createFiles(t *testing.T, directory string, files []sql.File) {
	t.Helper()

	for _, file := range files {
		path := filepath.Join(directory, file.Path)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("failed creating file: %v", err)
		}

		_ = f.Close()
	}
}

func prefixPaths(directory string, files []sql.File) []sql.File {
	f := slices.Clone(files)
	for i := range f {
		f[i].Path = filepath.Join(directory, f[i].Path)
	}

	return f
}
