BINARY      := ifscope
PKG         := github.com/bgrewell/ifscope
CMD         := ./cmd/ifscope
BIN_DIR     := bin
VERSION_PKG := $(PKG)/internal/version

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(VERSION_PKG).Version=$(VERSION) \
	-X $(VERSION_PKG).Commit=$(COMMIT) \
	-X $(VERSION_PKG).Date=$(DATE)

.PHONY: all build install test vet fmt fmt-check lint check tidy clean

all: check build

build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

install:
	go install -trimpath -ldflags "$(LDFLAGS)" $(CMD)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed; skipping"

check: fmt-check vet lint test

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)
