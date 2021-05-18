PROJECT_NAME := "simple-db-exporter"
PKG := "github.com/bisegni/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GO_ROOT := ${GOROOT}
VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

.PHONY: all dep build clean

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...
	@go mod download

build: dep ## Build the binary file
	echo "Build vesion: $(VERSION) build time: $(BUILDTIME)"
	@env GOOS=darwin go build $(GOFLAGS) -o simple-db-exporter-darwin -i -v $(PKG)
	@chmod a+x simple-db-exporter-darwin
	@tar cvzf simple-db-exporter-darwin.tar.gz simple-db-exporter-darwin

	@env GOOS=linux go build $(GOFLAGS) -o simple-db-exporter-linux -i -v $(PKG)
	@chmod a+x simple-db-exporter-linux
	@tar cvzf simple-db-exporter-linux.tar.gz simple-db-exporter-linux
	
	@env GOOS=windows go build $(GOFLAGS) -o simple-db-exporter-windows.exe -i -v $(PKG)
	@tar cvzf simple-db-exporter-windows.exe.tar.gz simple-db-exporter-windows.exe

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
