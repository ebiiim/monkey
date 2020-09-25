all: test build

test:
	go test -race -cover ./...

build: build-repl

build-repl:
	go build "-ldflags=-s -w" -trimpath -o monkey cmd/repl/main.go
