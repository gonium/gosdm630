TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

PWD := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BIN := $(PWD)/bin
# BUILD := GO111MODULE=on GOBIN=$(BIN) go install -v -ldflags '-X "main.version=${VERSION}" -X "main.commit=${SHA}" -X "main.date=${BUILD_DATE}"' ./...

all: build

build: assets binaries

binaries:
	@echo "Building for host platform"
	GO111MODULE=on GOBIN=$(BIN) go install -v -ldflags '-X "sdm630.Version=${VERSION}" -X "sdm630.Commit=${SHA}"' ./...
	@echo "Created binaries:"
	@ls -1 bin

assets:
	@echo "Generating embedded assets"
	GO111MODULE=on go generate ./...

release: test clean assets
	./build.sh

test:
	@echo "Running testsuite"
	GO111MODULE=on go test ./...

clean:
	rm -rf bin/ pkg/ *.zip

.PHONY: all build binaries assets release test clean
