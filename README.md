# Gator

Gator is a command-line RSS feed aggregator written in Go. It lets you register
users, follow RSS feeds, periodically scrape new posts, and browse the latest
posts from feeds you follow.

## Prerequisites

Before installing gator, make sure you have the following installed on your machine:

- [Go](https://go.dev/doc/install) (1.20 or later recommended)
- [PostgreSQL](https://www.postgresql.org/download/) (a running Postgres instance)

## Installation

Install the `gator` CLI using `go install`:

```bash
go install github.com/OrationKnight21/gator@latest
```

This downloads, compiles, and installs the `gator` binary into your
`$GOPATH/bin` (or `$GOBIN`) directory. Make sure that directory is in your
`PATH` so you can run `gator` from anywhere.

## Configuration

Gator reads its configuration from a JSON file located at `~/.gatorconfig.json`.

Create the file with the following structure:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

- `db_url`: the connection string for your Postgres database.
- `current_user_name`: gets set automatically when you log in; you can leave it blank initially.

Make sure your Postgres database exists and the schema migrations have been run
before using gator.

## Usage

Once installed and configured, you can run commands like this:

```bash
gator <command> [args...]
```

### Available commands

- `register <name>` — create a new user and log in as them.
- `login <name>` — log in as an existing user.
- `users` — list all registered users.
- `reset` — reset the database (deletes all users/feeds/posts).
- `addfeed <name> <url>` — add a new RSS feed (requires login) and automatically follow it.
- `feeds` — list all feeds that have been added.
- `follow <url>` — follow an existing feed by URL (requires login).
- `following` — list the feeds the current user is following (requires login).
- `unfollow <url>` — unfollow a feed (requires login).
- `agg <time_between_reqs>` — start the aggregator, which continuously scrapes
  feeds for new posts at the given interval (e.g. `1m`, `30s`, `1h`).
- `browse [limit]` — browse posts from feeds you follow (requires login).
- `search <query>` — fuzzy search posts by title or description (requires login).

### Example workflow

```bash
gator register alice
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml 
gator agg 1m
gator browse 5
```
```

## to use "search", just enter the following after registering the user.
gator search "golang" 10

## Notes

Go programs are statically compiled binaries. After running `go build` or
`go install`, you can run the resulting `gator` binary directly without needing
the Go toolchain installed on the machine that runs it. `go run .` is only for
local development.