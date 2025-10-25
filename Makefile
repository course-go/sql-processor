.PHONY: build
build:
	go build -o bin/sql-processor cmd/sql-processor/main.go

.PHONY: test
test:
	go test -cover -race -timeout 30s ./...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run
