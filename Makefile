PWD := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BIN := $(PWD)/bin
BUILD := env GO111MODULE=on GOBIN=$(BIN) go install ./...
GOPATH := $(shell go env GOPATH)

all: build

build: assets binaries

binaries:
	@echo "Building for host platform"
	$(BUILD)
	@echo "Created binaries:"
	@ls -1 bin

assets:
	./hash.sh
	@echo "Generating embedded assets"
	$(GOPATH)/bin/embed http.go

release: test clean assets
	./build.sh

test:
	@echo "Running testsuite"
	env GO111MODULE=on go test github.com/gonium/gosdm/...

clean:
	rm -rf bin/ pkg/ *.zip

dep:
	@echo "Installing embed tool"
	go get github.com/aprice/embed
	go install github.com/aprice/embed/cmd/embed

.PHONY: all build binaries assets release test clean dep
