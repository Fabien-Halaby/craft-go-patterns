.PHONY: test lint build run-store run-config clean

test:
	go test -race -cover ./...

lint:
	golangci-lint run

build:
	go build ./...

run-store:
	cd examples/singleton/store && go run main.go

run-config:
	cd examples/singleton/config && go run main.go

clean:
	go clean -cache

ci: lint test build