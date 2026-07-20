# Gator

[![Test](https://github.com/rvanhoepen/gator/actions/workflows/test.yml/badge.svg)](https://github.com/rvanhoepen/gator/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/rvanhoepen/gator/branch/main/graph/badge.svg)](https://codecov.io/gh/rvanhoepen/gator)

Gator is an RSS feed aggregator built in Go.

It is a command-line tool for collecting RSS feeds, storing posts in PostgreSQL, and reading summaries from the terminal. The name comes from aggreGATOR.

## Features

Gator is designed to help users keep up with blogs, news sites, podcasts, and other sources that publish RSS feeds.

- Add RSS feeds from across the internet
- Store collected feed posts in a PostgreSQL database
- Follow RSS feeds added by other users
- Unfollow feeds you no longer want to track
- Browse summaries of aggregated posts in the terminal

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
- `sqlc` if you plan to change SQL queries or schema and regenerate database code

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

## Configuration

Gator reads configuration from `~/.gatorconfig.json`.

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Fields:

| Field | Description |
| --- | --- |
| `db_url` | PostgreSQL connection string used by Gator |
| `current_user_name` | The currently logged-in user; managed by `register`, `login`, and `reset` |

You usually only need to edit `db_url` manually. Gator updates `current_user_name` when you register, log in, or reset users.

Run database migrations:

```sh
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

Run the project:

```sh
go run . <command> [args...]
```

## Commands

| Command | Description |
| --- | --- |
| `register <name>` | Create a user and log in as them |
| `login <name>` | Set the current user |
| `users` | List users |
| `reset` | Delete all users and clear the current user |
| `addfeed <name> <url>` | Add a feed and follow it as the current user |
| `feeds` | List all feeds |
| `follow <url>` | Follow an existing feed |
| `following` | List feeds followed by the current user |
| `unfollow <url>` | Unfollow a feed |
| `agg <duration>` | Continuously fetch feeds, e.g. `1m`, `30s` |
| `browse [limit]` | Show recent posts from followed feeds |

### User Commands

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

### Feed Commands

Add a feed and follow it as the current user:

```sh
go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

List all feeds:

```sh
go run . feeds
```

Follow an existing feed:

```sh
go run . follow "https://news.ycombinator.com/rss"
```

List feeds you follow:

```sh
go run . following
```

Unfollow a feed:

```sh
go run . unfollow "https://news.ycombinator.com/rss"
```

### Aggregation And Browsing

Start the feed aggregator:

```sh
go run . agg 1m
```

`agg` runs continuously and fetches one feed each interval. Durations use Go duration syntax, such as `30s`, `1m`, or `1h`. Stop it with `Ctrl+C`.

Browse recent posts from feeds you follow:

```sh
go run . browse 10
```

If no limit is provided, `browse` shows 2 posts.

## Example Workflow

```sh
go run . register ricardo
go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"
go run . agg 1m
go run . browse 10
```

## Why RSS?

RSS lets websites publish updates in a standard format. Instead of checking many sites manually, Gator can collect updates from all of them and show new posts in one place.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
