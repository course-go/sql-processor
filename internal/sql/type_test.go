package sql_test

import (
	"errors"
	"testing"

	"github.com/course-go/sql-processor/internal/sql"
)

func TestParseType(t *testing.T) {
	t.Parallel()

	t.Run("PostgresType", func(t *testing.T) {
		t.Parallel()

		sqlType, err := sql.ParseType("postgres")
		if err != nil {
			t.Fatalf("failed parsing valid type: %v", err)
		}

		expectedType := sql.PostgresType
		if sqlType != expectedType {
			t.Fatalf("types do not match: expect = %v, got = %v", expectedType, sqlType)
		}
	})

	t.Run("UnknownType", func(t *testing.T) {
		t.Parallel()

		_, err := sql.ParseType("unknown")
		if !errors.Is(err, sql.ErrUnknownType) {
			t.Fatal("expected unknown type error")
		}
	})
}
