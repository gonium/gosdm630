.PHONY: default clean checks test build assets binaries publish-images test-release

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

default: clean checks test build

clean:
	rm -rf bin/ pkg/ *.zip

checks: assets
	golangci-lint run

test:
	@echo "Running testsuite"
	GO111MODULE=on go test ./...

build: assets binaries

assets:
	@echo "Generating embedded assets"
	GO111MODULE=on go generate ./...

binaries:
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/gonium/gosdm630.Version=${VERSION}" -X "github.com/gonium/gosdm630.Commit=${SHA}"' ./cmd/sdm

publish-images:
	@echo Version: $(VERSION) $(BUILD_DATE)
	seihon publish -v "$(TAG_NAME)" -v "latest" --image-name andig/gosdm --base-runtime-image alpine --dry-run=false

test-release:
	goreleaser --snapshot --skip-publish --rm-dist
