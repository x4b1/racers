default: test lint

test:
	go test -race ./...

lint:
	staticcheck ./...

migration-create:
	migrate create -ext sql -dir storage/postgres/migrations -seq $(name)
