.PHONY: build build-migrate run test lint migrate-up migrate-down migrate-version

build:
	go build -o bin/server ./cmd/server

build-migrate:
	go build -o bin/migrate ./cmd/migrate

run:
	go run ./cmd/server

test:
	go test -v ./...

lint:
	golangci-lint run

migrate-up:
	./bin/migrate up

migrate-down:
	./bin/migrate down

migrate-version:
	./bin/migrate version
