# SQL Processor

SQL Processor is a command line application.
This project serves as a template for the second homework assignment.
To learn more about the homework assignments in general, visit
the [homework](https://github.com/course-go/homework) repository.

## Assignment

Throughout this homework assignment you will create SQL processor command line tool.
The tool receives its parameters over command line arguments and listens for SQL
file changes on local file system. These SQL files then get parsed and their SQL
statements get processed. After processing, they get exported to a bunch of exporters.

![SQL Processor diagram](assets/sql-processor.svg)

In the repository, you will find skeleton of the application. First of all, look
around. Some of the things you will need are already implemented.
You are free to change any code as long as all the tests work.
However, you should be able to finish the assignment just by implementing
the `Run` functions and some of the `New` constructors.
All the places that you should implement are denoted by a `TODO`.

### Specification

The basic syntax is as follows:

```shell
sql-processor [DIRECTORY DIRECTIVE]
```

where **directory directive** represents a pair of directory path and SQL dialect
separated using colon.
If there are no directives provided or some of the directives are invalid, the
application will exit with an error.

The example usage is as follows:

```shell
sql-processor ./sql/files:postgres /var/local/db/mysql:mysql
```

The application will then read all the specified directives, parse them and will
watch for newly created files in the given directories. For observing the file
system changes, the [fsnotify](https://github.com/fsnotify/fsnotify) library will
be used. When such newly created file appear, the application processes the file's
data and exports it.

You are safe to assume, that all the files created in the observed directories are
valid SQL files that may only contain single line comments. All SQL statements
in these files always start on a new line and always end with a semicolon.
SQL statement can however span multiple lines.

For example, this is a valid SQL file with 2 SQL statements:

```sql
-- Valid comment
SELECT u.name, o.total_amount, o.order_date
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.status = 'completed';

UPDATE users SET last_login = NOW() WHERE id = 123;
```

While this file is invalid because of the 2 following reasons:

```sql
-- Valid comment
SELECT u.name, o.total_amount, o.order_date
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.status = 'completed'; -- Nothing can follow a semicolon.

-- No multiple statements on a single line.
UPDATE users SET last_login = NOW() WHERE id = 123; SELECT * FROM users;
```

Your application may encountered errors during its runtime like not being
able to read a file etc. As we don't want to kill the app because of a single file,
most of the error processing will be done using logging. The components must
use the [slog](https://pkg.go.dev/log/slog) logger to log such errors
on appropriate level.

If you are unsure about some behaviour take the tests as the source of truth.

## Requirements

The CLI application has to support all the functionality previously
described. Stress is also given on writing idiomatic code and properly
handling resources and errors.

The application must run concurrently while not experiencing any deadlocks or
data races during its execution.

## Motivation

The main goal of this homework is to practice writing concurrent code in Go
using some of its key concepts like channels or contexts.

## Packages

Some of the Go packages worth looking into include:

- [fsnotify](https://github.com/fsnotify/fsnotify) for observing the
  file system changes
- [os](https://pkg.go.dev/os) and [filepath](https://pkg.go.dev/filepath) for
  interaction with file system
- [io](https://pkg.go.dev/io) and [bufio](https://pkg.go.dev/bufio) for interacting
  with input and outputs
- [sync](https://pkg.go.dev/sync) for synchronization
- [slog](https://pkg.go.dev/log/slog) for logging
