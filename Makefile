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
TAG=latest
LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"

.PHONY: check-gofmt
check-gofmt:
	diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: gofmt
gofmt:
	gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: ui
ui:
	cd ui && npm install && npm run bower install && npm run grunt build

.PHONY: release
release: image
	git tag -a $(TAG) -m $(TAG)
	git push origin $(TAG)
	goreleaser --rm-dist
	docker push docker.io/feedhenry/mcp-standalone:$(TAG)

build_cli:
	go build -o mcp ./cmd/mcp-cli

build: test-unit
	export GOOS=linux && go build ${LDFLAGS} ./cmd/mcp-api

image: build
	mkdir -p tmp
	cp ./mcp-api tmp
	cp artifacts/Dockerfile tmp
	cd tmp && docker build -t docker.io/feedhenry/mcp-standalone:$(TAG) .
	rm -rf tmp

run_server:
	@echo Running Server
	time go build ${LDFLAGS} ./cmd/mcp-api
	oc login -u developer -panything
	oc new-project $(NAMESPACE) | true
	oc create -f artifacts/openshift/sa.local.json -n  $(NAMESPACE) | true
	oc policy add-role-to-user edit system:serviceaccount:$(NAMESPACE):mcp-standalone -n  $(NAMESPACE) | true
	oc sa get-token mcp-standalone -n  $(NAMESPACE) > token
	./mcp-api -namespace=$(NAMESPACE) -k8-host=$(OSCP) -satoken-path=./token -log-level=debug -insecure=true

.PHONY: test
test: test-unit

.PHONY: setup
setup:
	@go get github.com/kisielk/errcheck

.PHONY: check
check:
	@echo Running checks:
	@echo errcheck
	@errcheck -ignoretests $$(go list ./...)
	@echo go vet
	@go vet ./...
	@echo go fmt
	diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

test-unit:
	@echo Running tests:
	go test -v -race -cover $(UNIT_TEST_FLAGS) \
	  $(addprefix $(PKG)/,$(TEST_DIRS))


apbs:
## Evaluate the presence of the TAG, to avoid evaluation of the nested shell script, during the read phase of make
    ifdef TAG
	@echo "Preparing $(TAG)"
        ifeq ($(shell git ls-files -m | wc -l),0)
			@echo "Doing the releae of the FeedHenry MCP APBs"
			cp artifacts//openshift/template.json cmd/android-apb/roles/provision-android-app/templates
			cp artifacts/openshift/template.json cmd/cordova-apb/roles/provision-cordova-apb/templates
			cp artifacts/openshift/template.json cmd/ios-apb/roles/provision-ios-apb/templates
			git commit -m "[make apbs script] updating Openshift template for APBs" cmd/
			cd cmd/android-apb && make build_and_push TAG=$(TAG)
			cd cmd/ios-apb && make build_and_push TAG=$(TAG)
			cd cmd/cordova-apb && make build_and_push TAG=$(TAG)
        else
	        $(error Aborting release process, since local files are modified)
        endif
    else
		$(error No VERSION defined!)
    endif

clean:
	./installer/clean.sh
