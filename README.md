# go-chi-web-api

Template to start a REST API with Go, Chi, Postgres, Pgx, Sqlc and Docker

## Requirements

- Go
- Git
- Docker

## Getting started

Click on the `Use this template` button to create a repository and then download.

Suppose I named the project as `mywebapi`, I would clone the repo with:

```sh
git clone git@github.com/opchaves/mywebapi.git
```

Install the needed tools and dependencies

```sh
cd mywebapi
make install-tools
make tidy
```

Create the env file

```sh
cp .env.example .env
```

Running the server in dev mode

```sh
# with live reload (air)
make dev
# or without live reload
make run
```

Building and running the binary

```sh
make build
make start
```
