# Kommonei

You are the king of your money. You rule your money, or your money will rule you.

Application to help you and family manage your personal finance.

## Tech Stack

- Backend: Go, Chi, Sqlc, Pgx, Docker
- Database: Postgres
- Frontend: React, Vite, Shadcn, Tailwind

## Requirements

- Go v1.22+
- Node v20+
- Git
- Docker

## Getting started

Clone the repository

```sh
git clone git@github.com:opchaves/kommonei.git
```

Create the env file

```sh
cp .env.example .env
```

Install the needed tools and dependencies

```sh
cd kommonei
make install-tools
make tidy
```

Running the server in dev mode

```sh
# with live reload (air)
make dev
# or without live reload
make run
```

Landing page: [localhost:8080](http://localhost:8080)

API: [localhost:8080/api](http://localhost:8080/api)

Building and running the binary

```sh
make build
make start
```

## Links & Credits

- [Valkyrie: A Discord clone using React and Go.](https://github.com/sentrionic/Valkyrie)
- [Go RESTful API Boilerplate with JWT Authentication backed by PostgreSQL](https://github.com/dhax/go-base)
- [Go Rest API starter kit / Golang API boilerplate base on Chi framework](https://github.com/qreasio/go-starter-kit)
- [A Golang and HTMX Todo application](https://github.com/paganotoni/todox)
- [Pico.css examples](https://github.com/picocss/examples)
- [Pico V2: A pure HTML example, without dependencies](https://codesandbox.io/s/github/picocss/examples/tree/master/v2-html)
- [Cache-Control on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control)
- [How to setup a PostgreSQL database with Docker Compose](https://blog.cadumagalhaes.dev/how-to-setup-a-postgresql-database-with-docker-compose)
