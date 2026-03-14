BINARY = earnings-cal-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X github.com/oinkywaddles/earnings-cal-cli/cmd.version=$(VERSION)"

.PHONY: build clean test lint

build:
	go build $(LDFLAGS) -o $(BINARY) .

clean:
	rm -f $(BINARY)

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install $(LDFLAGS) .
