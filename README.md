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
- [Landing page template using Shadcn, React, Typescript and Tailwind](https://github.com/leoMirandaa/shadcn-landing-page)
- [Cache-Control on MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control)
- [How to setup a PostgreSQL database with Docker Compose](https://blog.cadumagalhaes.dev/how-to-setup-a-postgresql-database-with-docker-compose)
- [ToDo App - Chi router, HTMx](https://github.com/Ujstor/todo-go-htmx)
- [How to load and export variables from an .env file in Makefile?](https://stackoverflow.com/a/69902063)
- [Create the smallest and secured golang docker image based on scratch](https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324)
- [Get the short git version hash](https://stackoverflow.com/questions/5694389/get-the-short-git-version-hash)

## TODO

- [ ] Is a nginx config needed when running on Digital Ocean?

  - App will be served from Go
  - If app were server from nginx. [Setup docker to run React app](https://gist.github.com/przbadu/929fc2b0d5d4cd78a5efe76d37f891b6)

- [ ] Deploy to Digital Ocean using its container registry (DOCR)

  [DigitalOcean Container Registry](https://docs.digitalocean.com/products/container-registry/)
