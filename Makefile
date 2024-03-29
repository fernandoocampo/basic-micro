APP_VERSION?=0.7.2
IMAGE?=pets:$(APP_VERSION)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
SRC_FOLDER=cmd/petsd
BINARY_NAME=pets
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
BINARY_DARWIN=$(BINARY_NAME)-amd64-darwin
COMMIT_HASH?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION?=$(shell git describe --dirty --tags --always)
LDFLAGS?="-X github.com/fernandoocampo/basic-micro/internal/service.Version=${VERSION} -X github.com/fernandoocampo/basic-micro/internal/service.CommitHash=${COMMIT_HASH} -X github.com/fernandoocampo/basic-micro/internal/service.BuildDate=${BUILD_DATE} -s -w"

.PHONY: clean
clean:  ## clean binaries
	$(GOCLEAN)
	rm bin/$(BINARY_DARWIN)
	rm bin/$(BINARY_UNIX)

.PHONY: build
build: ## Build binary for mac
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN} ./${SRC_FOLDER}/main.go

.PHONY: build-linux
build-linux: ## Build binary for Linux
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_UNIX} ./${SRC_FOLDER}/main.go

.PHONY: run-local
run-local: ## run local
	go run -ldflags ${LDFLAGS} cmd/petsd/main.go

.PHONY: run-docker-local
run-docker-local: ## run project local
	docker-compose up --build

.PHONY: clean-docker-local
clean-docker-local: ## clean docker-compsoe
	docker-compose down

.PHONY: test
test: ## run unit tests
	${GOCMD} test -race ./...

.PHONY: vet
vet: ## run unit tests
	${GOCMD} vet ./...

.PHONY: print-image-name
print-image-name: ## print current app version
	echo ${IMAGE}