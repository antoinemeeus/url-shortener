# Url-Shortener

An URL shortener service, basically an overglorified hashmap of urls, written in Golang.

## Requirement

Go version: 1.14

## Start App locally

The URL shortenenr can use [Redis](https://redis.io/) or [Postgres](https://www.postgresql.org/) as database.

You can choose between Redis or Postgres by modifying the `DB_ENGINE` variable in the `.env` file.
By default Postgres is used: `DB_ENGINE=postgres`

See `docker-compose.yaml` definition for running and configuring databases images locally.

You can start the service locally by running:
```bash
make start
```

App will be available on `localhost:8000`. 

You can change the port with `PORT=8000` in `.env` file.

