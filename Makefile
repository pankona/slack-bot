
all: test lint build

build:
	go build ./...

test:
	go test ./...

lint:
	golint ./...

clean:
	rm -f $(CURDIR)/slack-bot
	rm -f $(CURDIR)/cmd/slack-bot/slack-bot
