APP_VERSION?=0.7.3
IMAGE?=pets:$(APP_VERSION)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
CONTAINERCMD?=docker
SRC_FOLDER=cmd/petsd
BINARY_NAME=pets
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
BINARY_DARWIN=$(BINARY_NAME)-amd64-darwin
COMMIT_HASH?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION?=$(shell git describe --dirty --tags --always)
LDFLAGS?="-X github.com/fernandoocampo/basic-micro/internal/application.Version=${VERSION} -X github.com/fernandoocampo/basic-micro/internal/application.CommitHash=${COMMIT_HASH} -X github.com/fernandoocampo/basic-micro/internal/application.BuildDate=${BUILD_DATE} -s -w"
DOCKER_REPO?=fdocampo

.PHONY: clean
clean:  ## clean binaries
	$(GOCLEAN)
	rm bin/$(BINARY_DARWIN)
	rm bin/$(BINARY_UNIX)

.PHONY: build-linux
build-linux: ## Build binary for Linux
	${GOCMD} mod tidy
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_UNIX} ./${SRC_FOLDER}/main.go

.PHONY: build-mac
build-mac: ## Build binary for mac
	${GOCMD} mod tidy
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN} ./${SRC_FOLDER}/main.go

.PHONY: build-image
build-image: ## build container image
	${CONTAINERCMD} build \
	--build-arg appVersion=${VERSION} \
	--build-arg buildDate=${BUILD_DATE} \
	--build-arg commitHash=${COMMIT_HASH} \
    -f Dockerfile \
    -t basic-micro:${VERSION} .

.PHONY: run-container-local
run-container-local: ## Run new container local
	${CONTAINERCMD} run --rm -it -p 8080:8080 \
	basic-micro:${VERSION}

.PHONY: login-docker-hub
login-docker-hub:
	${CONTAINERCMD} login docker.io

.PHONY: docker-image-push
docker-image-push:
	${CONTAINERCMD} tag basic-micro:${VERSION} ${DOCKER_REPO}/basic-micro:${VERSION}
	${CONTAINERCMD} push ${DOCKER_REPO}/basic-micro:${VERSION}

# docker tag local-image:tagname new-repo:tagname
# docker push new-repo:tagname

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

.PHONY: show-image-history
show-image-history: ## check layers of an image and their size
	${CONTAINERCMD} history basic-micro:${VERSION}