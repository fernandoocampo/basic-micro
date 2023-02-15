# Builder image go
FROM golang:1.19.4 AS builder

ARG appVersion
ARG commitHash

ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash

# Build pets binary with Go
ENV GOPATH /opt/go

RUN mkdir -p /pets
WORKDIR /pets
COPY . /pets
RUN go mod tidy
RUN make build-linux

# Runnable image
FROM alpine
ARG appVersion
ARG commitHash
ENV VERSION=$appVersion
ENV COMMIT_HASH=$commitHash
COPY --from=builder /pets/bin/pets-amd64-linux /bin/pets-service
RUN ls /bin/pets-service
WORKDIR /bin
ENTRYPOINT [ "./pets-service" ]