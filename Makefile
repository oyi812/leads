DIRS=bin

all: test build

build:
	@go build -o bin/leads ./cmd/leads
test:
	@go test -v ./internal/...

$(info $(shell mkdir -p $(DIRS)))
