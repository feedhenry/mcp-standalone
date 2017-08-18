PKG     = github.com/feedhenry/mobile-server
TOP_SRC_DIRS   = pkg
TEST_DIRS     ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
                   -exec dirname {} \\; | sort | uniq")

build:
	export GOOS=linux && go build ./cmd/mobile-server 


image: build
	mkdir -p tmp
	cp ./mobile-server tmp
	cp artifacts/Dockerfile tmp
	docker build -t feedhenry/mobile-server:latest tmp
	docker tag feedhenry/mobile-server:latest feedhenry/mobile-server:latest
	rm -rf tmp

test: build test-unit

test-unit: build
	@echo Running tests:
	go test -cover $(UNIT_TEST_FLAGS) \
	  $(addprefix $(PKG)/,$(TEST_DIRS))
	