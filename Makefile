all: test

test:
	go test -race -cover ./...
