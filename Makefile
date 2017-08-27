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

.PHONY: web
web:
	gem install compass
	cd web && npm install && ./node_modules/.bin/bower install && grunt build

build: web test-unit
	export GOOS=linux && go build ./cmd/mcp-standalone

image: build
	mkdir -p tmp
	mkdir -p tmp/web/dist
	cp ./mcp-standalone tmp
	cp artifacts/Dockerfile tmp
	cp -R web/dist tmp/web/dist
	cd tmp && docker build -t feedhenry/mcp-standalone:latest .
	rm -rf tmp

run:
	@echo Running Server
	go install ./cmd/mcp-standalone
	oc new-project test | true
	oc create sa mobile-server | true
	oc sa get-token mobile-server >> token
	mcp-standalone -namespace=test -k8-host=https://192.168.37.1:8443 -satoken-path=./token


test: test-unit

test-unit:
	@echo Running tests:
	go test -cover $(UNIT_TEST_FLAGS) \
	  $(addprefix $(PKG)/,$(TEST_DIRS))
	
