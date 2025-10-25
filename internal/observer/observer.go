package observer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/course-go/sql-processor/internal/sql"
	"github.com/fsnotify/fsnotify"
)

const directivePartCount = 2

var (
	ErrNoDirectoryDirectivesProvided = errors.New("no directory paths provided")
	ErrInvalidDirectoryDirective     = errors.New("directory directive has invalid format")
)

// Observer observes given filesystem directories for new files.
// When it notices such file it creates a [sql.File] and passes it for processing.
type Observer struct {
	logger         *slog.Logger
	watcher        *fsnotify.Watcher
	fileCh         chan<- sql.File
	directoryTypes map[string]sql.Type
}

// New creates a new [Observer].
//
// The directives parameter represents a directory directives in the "[directory]:[sql.Type]" format.
// For example, the following is a valid directory directive: "/var/sql/postgres:postgres".
func New(logger *slog.Logger, directives []string, fileCh chan<- sql.File) (o Observer, err error) {
	// TODO: Implement.

	return Observer{}, nil
}

// Run starts the [Observer].
func (o *Observer) Run(ctx context.Context) {
	// TODO: Implement.
}

// Close closes the [Observer].
func (o *Observer) Close() error {
	// TODO: Implement.

	return nil
}
