
all: test lint build

build:
	go build ./...

test:
	go test ./...

lint:
	golint ./...
