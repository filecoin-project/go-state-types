GO_BIN ?= go
GOLINT ?= golangci-lint

all: build lint test tidy
.PHONY: all

build:
	$(GO_BIN) build ./...
.PHONY: build

test:
	$(GO_BIN) test ./...
.PHONY: test

test-coverage:
	$(GO_BIN) test -coverprofile=coverage.out ./...
.PHONY: test-coverage

tidy:
	$(GO_BIN) mod tidy
.PHONY: tidy

gen:
	$(GO_BIN) run ./gen/gen.go
	$(GO_BIN) run ./manifest/gen/gen.go
	$(GO_BIN) run ./proof/gen/gen.go
	$(GO_BIN) run ./batch/gen/gen.go
	$(GO_BIN) run ./builtin/v8/gen/gen.go
	$(GO_BIN) run ./builtin/v9/gen/gen.go
	$(GO_BIN) run ./builtin/v10/gen/gen.go
	$(GO_BIN) run ./builtin/v11/gen/gen.go
	$(GO_BIN) run ./builtin/v12/gen/gen.go
	$(GO_BIN) run ./builtin/v13/gen/gen.go
	$(GO_BIN) run ./builtin/v14/gen/gen.go
	$(GO_BIN) run ./builtin/v15/gen/gen.go
	$(GO_BIN) run ./builtin/v16/gen/gen.go
	$(GO_BIN) run ./builtin/v17/gen/gen.go	
.PHONY: gen

lint:
	$(GOLINT) run ./...
.PHONY: lint
