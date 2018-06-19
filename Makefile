PWD := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BIN := $(PWD)/bin
BUILD := GOBIN=$(BIN) go install ./...

all: build

build: assets binaries

binaries:
	@echo "Building for host platform"
	@$(BUILD)
	@echo "Created binaries:"
	@ls -1 bin

assets:
	@echo "Generating embedded assets"
	@$(GOPATH)/bin/embed http.go

release-build: test clean assets
	@echo "Building binaries..."
	@echo "... for Linux/32bit"
	@GOOS=linux GOARCH=386 $(BUILD)
	@echo "... for Linux/64bit"
	@GOOS=linux GOARCH=amd64 $(BUILD)
	@echo "... for Raspberry Pi/Linux"
	@GOOS=linux GOARCH=arm GOARM=5 $(BUILD)
	@echo "... for Mac OS/64bit"
	@GOOS=darwin GOARCH=amd64 $(BUILD)
	@echo "... for Windows/32bit"
	@GOOS=windows GOARCH=386 $(BUILD)
	@echo "... for Windows/64bit"
	@GOOS=windows GOARCH=amd64 $(BUILD)
	@echo
	@echo "Created binaries:"
	@ls -1 bin

package:
	@echo "Starting packaging"
	@echo "... for Linux"
	@zip sdm630-linux-386 bin/*-linux-386*
	@zip sdm630-linux-amd64 bin/*-linux-amd64
	@zip sdm630-linux-arm bin/*-linux-arm*
	@echo "... for Mac OS"
	@zip sdm630-darwin-amd64 bin/*-darwin-amd64
	@echo "... for Windows"
	@zip sdm630-windows-386 bin/*-windows-386*

release: release-build package

test:
	@echo "Running testsuite"
	@go test

clean:
	@rm -rf bin/ pkg/ *.zip

dep:
	@echo "Installing embed tool"
	@go get -u github.com/aprice/embed/cmd/embed

.PHONY: all build binaries assets release-build package release test clean dep
