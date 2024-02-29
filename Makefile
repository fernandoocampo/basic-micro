APP_VERSION?=0.5.1
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
LDFLAGS?="-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -s -w"

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

.PHONY: run-app
run-app: ## run app
	go run cmd/petsd/main.go

.PHONY: run-local
run-local: ## run project local
	docker-compose up --build

.PHONY: clean-local
clean-local: ## clean docker-compsoe
	docker-compose down

.PHONY: test
test: ## run unit tests
	${GOCMD} test -race ./...

.PHONY: print-image-name
print-image-name: ## print current app version
	echo ${IMAGE}