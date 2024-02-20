# pets service

This is a service that provides an API for pets administration.

## How to build?

from project folder run below commands, it will output binaries in `./bin/` folder

* just build with current operating system.
```sh
make build
```

* build for a linux distro operating system.
```sh
make build-linux
```

## How to run a test environment quickly?

1. make sure you have docker-compose installed.
2. run the docker compose.
```sh
docker-compose up --build
```

or run this shortcut:

```sh
make run-local
```

3. once you are done using the environment follow these steps.

    * ctrl + c
    * make clean-local

## How to test?

from project folder run the following command

```sh
make test
```

or 

```sh
go test -race ./...
```

## Adding open telemetry

you can follow these [instructions](https://opentelemetry.io/docs/instrumentation/go/getting-started/)

```sh
go get "go.opentelemetry.io/otel" \
  "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric" \
  "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" \
  "go.opentelemetry.io/otel/propagation" \
  "go.opentelemetry.io/otel/sdk/metric" \
  "go.opentelemetry.io/otel/sdk/resource" \
  "go.opentelemetry.io/otel/sdk/trace" \
  "go.opentelemetry.io/otel/semconv/v1.21.0" \
  "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
```
