PKG     = github.com/feedhenry/mcp-standalone
TOP_SRC_DIRS   = pkg
TEST_DIRS     ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
                   -exec dirname {} \\; | sort | uniq")
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
SHELL = /bin/bash
#CHANGE this if using a different url for openshift
OSCP = https://192.168.37.1:8443

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

run_server:
	@echo Running Server
	time go install ./cmd/mcp-standalone
	oc new-project mcp-standalone | true
	oc create -f install/openshift/sa.local.json -n  mcp-standalone | true
	oc sa get-token mcp-standalone -n  mcp-standalone >> token
	mcp-standalone -namespace=mcp-standalone -k8-host=$(OSCP) -satoken-path=./token


test: test-unit

test-unit:
	@echo Running tests:
	go test -cover $(UNIT_TEST_FLAGS) \
	  $(addprefix $(PKG)/,$(TEST_DIRS))
	
