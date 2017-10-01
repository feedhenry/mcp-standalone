PKG     = github.com/feedhenry/mcp-standalone
TOP_SRC_DIRS   = pkg
TEST_DIRS     ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
                   -exec dirname {} \\; | sort | uniq")
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
SHELL = /bin/bash
#CHANGE this if using a different url for openshift
OSCP = https://192.168.37.1:8443
NAMESPACE =project2

.PHONY: check-gofmt
check-gofmt:
	diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: gofmt
gofmt:
	gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: ui
ui:
	cd ui && npm install && npm run bower install && npm run grunt build

build_cli:
	go build -o mcp ./cmd/mcp-cli

build: test-unit
	export GOOS=linux && go build ./cmd/mcp-api

image: build
	mkdir -p tmp
	cp ./mcp-api tmp
	cp artifacts/Dockerfile tmp
	cd tmp && docker build -t feedhenry/mcp-standalone:latest .
	rm -rf tmp

run_server:
	@echo Running Server
	time go build ./cmd/mcp-api
	oc login -u developer -panything
	oc new-project $(NAMESPACE) | true
	oc create -f install/openshift/sa.local.json -n  $(NAMESPACE) | true
	oc policy add-role-to-user edit system:serviceaccount:$(NAMESPACE):mcp-standalone -n  $(NAMESPACE) | true
	oc sa get-token mcp-standalone -n  $(NAMESPACE) > token
	./mcp-api -namespace=$(NAMESPACE) -k8-host=$(OSCP) -satoken-path=./token -log-level=debug -insecure=true


test: test-unit

test-unit:
	@echo Running tests:
	go test -v -race -cover $(UNIT_TEST_FLAGS) \
	  $(addprefix $(PKG)/,$(TEST_DIRS))

apbs:
	cp install/openshift/template.json cmd/android-apb/roles/provision-android-app/templates
	cp install/openshift/template.json cmd/cordova-apb/roles/provision-cordova-apb/templates	
	cp install/openshift/template.json cmd/ios-apb/roles/provision-ios-apb/templates
	cd cmd/android-apb && make build_and_push 		
	cd cmd/ios-apb && make build_and_push 		
	cd cmd/cordova-apb && make build_and_push

clean:
	./ui/clean.sh
