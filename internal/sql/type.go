package sql

import "errors"

var ErrUnknownType = errors.New("unknown sql type")

// Type represents SQL dialect.
type Type string

const (
	PostgresType Type = "postgres"
	MySQL        Type = "mysql"
	SQLite       Type = "sqlite"
)

func ParseType(input string) (t Type, err error) {
	switch Type(input) {
	case PostgresType:
		return PostgresType, nil
	case MySQL:
		return MySQL, nil
	case SQLite:
		return SQLite, nil
	default:
		return "", ErrUnknownType
	}
}
