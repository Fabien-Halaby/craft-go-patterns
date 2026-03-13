.PHONY: test lint build demo clean

test:
	go test -race -cover ./...

lint:
	golangci-lint run

build:
	go build ./...

demo:
	cd examples/demo && go run main.go

clean:
	go clean -cache

ci: lint test build