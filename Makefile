
all: build test

build:
	go build ./...

test:
	go test -v ./...

test-short:
	go test ./...

lint:
	go vet ./...

tidy:
	go mod tidy

clean:
	go clean ./...
