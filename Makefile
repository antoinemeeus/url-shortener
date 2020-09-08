PROJECT_NAME := "url-shortener"
PKG := "github.com/antoinemeeus/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep lint vet test test-coverage build clean

all:: help

help ::
	@grep -E '^[a-zA-Z_-]+\s*:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

start :: ## Run service
	@echo "  > Starting service"
	@docker-compose up -d
	@go run cmd/server/main.go

dep: ## Get the dependencies
	@go mod download

lint: ## Lint Golang files
	@golangci-lint run ./...

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unit tests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@go build -i -o build/server $(PKG)/cmd/server

clean: ## Remove previous build
	@rm -f ./build/server
