# pets service

This is a service that provides an API for pets administration.

## How to build?

from project folder run below commands, it will output binaries in `./bin/` folder

* just build with current operating system
```sh
make build
```

* build for a linux distro operating system
```sh
make build-linux
```

## How to run a test environment quickly?

1. make sure you have docker-compose installed.
2. run the docker compose.
```sh
docker-compose up --build
```

or run this shortcut

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
