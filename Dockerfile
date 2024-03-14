# Builder image go
FROM golang:1.22.1 AS builder

ARG appVersion
ARG commitHash

ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash

# Build pets binary with Go
ENV GOPATH /opt/go

RUN mkdir -p /pets
WORKDIR /pets
COPY . /pets
RUN go mod tidy && make build-linux

# Runnable image
FROM alpine:3.19
ARG appVersion
ARG commitHash
ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash
COPY --from=builder /pets/bin/pets-amd64-linux /bin/pets-service
RUN ls /bin/pets-service
WORKDIR /bin
ENTRYPOINT [ "./pets-service" ]