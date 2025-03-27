default: fmt lint build test

fmt:
	go fmt ./...

lint:
	golangci-lint run

build:
	go build ./...

test:
	go test ./...

benchmark:
	go test -bench=. -benchmem ./...
