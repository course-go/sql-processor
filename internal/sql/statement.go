package sql

// Statement represents SQL statement in SQL file.
type Statement struct {
	File    File
	Content string
	LineNum int
}
