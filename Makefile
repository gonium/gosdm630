.PHONY: default clean checks build binaries assets release test test-release

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

default: clean checks test build

clean:
	rm -rf bin/ pkg/ *.zip

checks:
	golangci-lint -e U1000 -e sunspecModelID run

build: assets binaries

binaries:
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "sdm630.Version=${VERSION}" -X "sdm630.Commit=${SHA}"' ./...

assets:
	@echo "Generating embedded assets"
	GO111MODULE=on go generate ./...

release: test clean assets
	./build.sh

test:
	@echo "Running testsuite"
	GO111MODULE=on go test ./...

test-release:
	goreleaser --snapshot --skip-publish --rm-dist
