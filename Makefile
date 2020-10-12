default: test lint

test:
	go test -race ./...

lint:
	golangci-lint run

migration-create:
	migrate create -ext sql -dir storage/postgres/migrations -seq $(name)
