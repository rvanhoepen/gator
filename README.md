# Gator 🐊

Gator is an RSS feed aggregator built in Go.

It is a command-line tool for collecting RSS feeds, storing posts in PostgreSQL, and reading summaries from the terminal. The name comes from aggreGATOR.

## What It Does

Gator is designed to help users keep up with blogs, news sites, podcasts, and other sources that publish RSS feeds.

Planned features include:

- Add RSS feeds from across the internet
- Store collected feed posts in a PostgreSQL database
- Follow RSS feeds added by other users
- Unfollow feeds you no longer want to track
- Browse summaries of aggregated posts in the terminal
- Open or copy links to the full original posts

## Tech Stack

- Go
- PostgreSQL
- RSS/XML feed parsing
- Command-line interface

## Requirements

To run Gator locally, you will need:

- Go installed
- PostgreSQL installed and running
- A PostgreSQL database for Gator
- `goose` for database migrations

## Status

This project is currently in early development.

## Getting Started

Clone the repository:

```sh
git clone <repository-url>
cd gator
```

Create a PostgreSQL database:

```sh
createdb gator
```

Create a config file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Update `db_url` to match your local PostgreSQL username, password, host, port, and database name.

Run database migrations:

```sh
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

Run the project:

```sh
go run . <command> [args...]
```

## Commands

Register a user and set them as the current user:

```sh
go run . register ricardo
```

Log in as an existing user:

```sh
go run . login ricardo
```

List all users:

```sh
go run . users
```

Delete all users and clear the current user:

```sh
go run . reset
```

## Why RSS?

RSS lets websites publish updates in a standard format. Instead of checking many sites manually, Gator can collect updates from all of them and show new posts in one place.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
