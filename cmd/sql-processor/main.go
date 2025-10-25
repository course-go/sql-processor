package main

import (
	"context"
	"fmt"
	"os"

	"github.com/course-go/sql-processor/internal/cmd"
	"github.com/course-go/sql-processor/internal/exporter"
	"github.com/course-go/sql-processor/internal/exporter/stdout"
)

func main() {
	exporters := []exporter.Exporter{stdout.NewExporter()}

	err := cmd.Run(context.Background(), os.Args, exporters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed running sql processor: %v", err)
		os.Exit(1)
	}
}
