BIN             = go-microservice-1
OUTPUT_DIR      = build
TMP_DIR        := .tmp
RELEASE_TIME   := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
RELEASE_VER    := $(shell git rev-parse --short HEAD)
LDFLAGS        := "-s -w -X main.version=$(RELEASE_VER)-$(RELEASE_TIME)"
DOCKER_IP       = $(shell docker info | grep -q moby && echo localhost || docker-machine ip)
NAME            = default
COVERMODE       = atomic

TEST_PACKAGES      := $(shell go list ./... | grep -v vendor | grep -v fakes | grep -v ftest)

.PHONY: help docs
.DEFAULT_GOAL := help

# if a .env.local file exists, use that instead
ifneq ("$(wildcard .env.local)","")
	FILE_FLAG := -e .env.local
endif

run: ## Run application (without building)
	go run *.go -d $(FILE_FLAG)

all: test build docker ## Test, build and docker image build

setup: installtools ## Install and setup tools


## Run Tests ##

# Under the hood, `go test -tags ...` also runs the "default" (unit) test case
# in addition to the specified tags
test: test/integration ## Perform both unit and integration tests

testv: testv/integration ## Perform both unit and integration tests (with verbose flags)

test/unit: ## Perform unit tests
	go test -cover $(TEST_PACKAGES)

testv/unit: ## Perform unit tests (with verbose flag)
	go test -v -cover $(TEST_PACKAGES)

test/integration: ## Perform integration tests
	go test -cover -tags integration $(TEST_PACKAGES)

testv/integration: ## Perform verbose integration tests
	go test -v -cover -tags integration $(TEST_PACKAGES)

test/race: ## Perform unit tests and enable the race detector
	go test -race -cover $(TEST_PACKAGES)

test/cover: ## Run all tests + open coverage report for all packages
	echo 'mode: $(COVERMODE)' > .coverage
	for PKG in $(TEST_PACKAGES); do \
		go test -coverprofile=.coverage.tmp -tags "integration" $$PKG; \
		grep -v -E '^mode:' .coverage.tmp >> .coverage; \
	done
	go tool cover -html=.coverage
	$(RM) .coverage .coverage.tmp

docs: ## Generate documentation (make sure you've ran `make installtools`)
	swag init

installtools: ## Install development related tools
	go get github.com/swaggo/swag
	go get github.com/maxbrunsfeld/counterfeiter

## Build ##

build: clean build/linux build/darwin ## Build for linux and darwin (save to OUTPUT_DIR/BIN)

build/linux: clean/linux ## Build for linux (save to OUTPUT_DIR/BIN)
	GOOS=linux go build -a -installsuffix cgo -ldflags $(LDFLAGS) -o $(OUTPUT_DIR)/$(BIN)-linux .

build/darwin: clean/darwin ## Build for darwin (save to OUTPUT_DIR/BIN)
	GOOS=darwin go build -a -installsuffix cgo -ldflags $(LDFLAGS) -o $(OUTPUT_DIR)/$(BIN)-darwin .

generate: ## Run generate for non-vendor packages only
	go list ./... | grep -v vendor | xargs go generate
	go fmt ./fakes/...

docker: build/linux ## Build local docker image
	docker build -t $(BIN):$(RELEASE_VER) -f Dockerfile .

docker/local: build/linux ## Bring up the service via docker-compose
	docker-compose build
	docker stop $(BIN); docker rm $(BIN); docker-compose up

jet: ## Run `jet steps`
	jet steps

clean: clean/darwin clean/linux ## Remove all build artifacts

clean/darwin: ## Remove darwin build artifacts
	$(RM) $(OUTPUT_DIR)/$(BIN)-darwin

clean/linux: ## Remove linux build artifacts
	$(RM) $(OUTPUT_DIR)/$(BIN)-linux

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'