PKG     = github.com/feedhenry/mcp-standalone
TOP_SRC_DIRS   = pkg
TEST_DIRS     ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
                   -exec dirname {} \\; | sort | uniq")
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
SHELL = /bin/bash

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor

.PHONY: check-gofmt
check-gofmt:
	diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: gofmt
gofmt:
	gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

build: test-unit
	export GOOS=linux && go build ./cmd/mcp-standalone


image: build
	mkdir -p tmp
	cp ./mobile-server tmp
	cp artifacts/Dockerfile tmp
	docker build -t feedhenry/mcp-standalone:latest tmp
	docker tag feedhenry/mcp-standalone:latest feedhenry/mcp-standalone:latest
	rm -rf tmp

test: test-unit

test-unit:
	@echo Running tests:
	go test -cover $(UNIT_TEST_FLAGS) \
	  $(addprefix $(PKG)/,$(TEST_DIRS))
	
