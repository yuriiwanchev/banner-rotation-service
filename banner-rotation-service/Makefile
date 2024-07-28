BIN := "./bin/banner-rotation-service"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/main.go

run: build
	$(BIN) -config ./configs/config.toml

test:
	go test -race -count 100 ./...
