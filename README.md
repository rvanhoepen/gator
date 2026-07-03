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

## Status

This project is currently in early development.

## Getting Started

Clone the repository:

```sh
git clone <repository-url>
cd gator
```

Initialize the Go module if it has not been created yet:

```sh
go mod init gator
```

Create a PostgreSQL database:

```sh
createdb gator
```

Run the project:

```sh
go run .
```

## Example Usage

The command interface is still being built, but Gator will support commands similar to:

```sh
gator addfeed "Tech Blog" "https://example.com/rss.xml"
gator follow "Tech Blog"
gator unfollow "Tech Blog"
gator browse
```

## Why RSS?

RSS lets websites publish updates in a standard format. Instead of checking many sites manually, Gator can collect updates from all of them and show new posts in one place.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
