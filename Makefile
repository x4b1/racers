default: test lint

test:
	go test -race ./...
lint:
	golangci-lint run
