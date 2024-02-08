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

Landing page: [localhost:8080](http://localhost:8080)

API: [localhost:8080](http://localhost:8080/api)

Building and running the binary

```sh
make build
make start
```

## Links

- [A Golang and HTMX Todo application](https://github.com/paganotoni/todox)
- [Pico.css examples](https://github.com/picocss/examples)
- [Pico V2: A pure HTML example, without dependencies](https://codesandbox.io/s/github/picocss/examples/tree/master/v2-html)
- [Cache-Control on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control)
